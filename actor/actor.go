package actor
import (

	"fmt"
	"sync/atomic"
	"reflect"
	"strconv"
	"sync"
	"queue"
)




/**
Actual actors kept within the system
Children are kept in a slice, also a map
This brings memory overhead, but must be a bit faster if it was to be kept
in a map only
 */
type DefaultActor struct {
	Name           string
	ChildrenArray  []*DefaultActor
	ChildrenMap    map[string]*DefaultActor
	Parent         *DefaultActor
	actorInterface ActorInterface
	messageQueue   *queue.RoundRobinQueue
	index          int     //todo is it going to be used
	ActorSystem    *ActorSystem

	stop           *uint32 //todo Find atomic operations for 8 bit or boolean
	pendingStop    *uint32

	justStarted    *uint32

	restart        *uint32

	stopped        *uint32

	startStopLock  *sync.Mutex
	selfLock       *sync.Mutex
}

func NewDefaultActor(name string, parent *DefaultActor) *DefaultActor {
	da := DefaultActor{}
	da.Name = name
	da.Parent = parent

	da.stop = new(uint32)
	*da.stop = 0

	da.pendingStop = new(uint32)
	*da.pendingStop = 0

	da.justStarted = new(uint32)
	*da.justStarted = 1

	da.restart = new(uint32)
	*da.restart = 0

	da.stopped = new(uint32)
	*da.stopped = 0

	da.ActorSystem = parent.ActorSystem

	da.messageQueue = queue.NewRoundRobinQueue()

	da.startStopLock = new(sync.Mutex)
	da.selfLock = new(sync.Mutex)

	return &da
}

//Set priority for message
func (selfPtr *DefaultActor) SetMessagePriority(typ reflect.Type, priority int) {
	selfPtr.messageQueue.SetGroup(getTypeNameFromType(typ), priority)
}

/**
This Block sets, resets and checks states
Uint32 type could be changed for something else later using atomic.Value
I am not sure which one is more efficient. A benchmark could be good
 */

func (selfPtr *DefaultActor) getMessageQueueCount() int {
	return selfPtr.messageQueue.GetTotalItemCount()
}

func (selfPtr *DefaultActor) atomicallySetPtr(ptr *uint32) {
	atomic.CompareAndSwapUint32(ptr, 0, 1)
}

func (selftPtr *DefaultActor) atomicallyResetPtr(ptr *uint32) {
	atomic.CompareAndSwapUint32(ptr, 1, 0)
}

func (selfPtr *DefaultActor) atomicallyCheckPtr(ptr *uint32) uint32 {
	return atomic.LoadUint32(ptr)
}

/**
Actor instance here will be a pointer as long as the user passes the pointer of the interface they implemented
//todo How to make them do it all the time to avoid mistakes in their behalf
 */
func recoverFunc(actor ActorInterface) {
	if r := recover(); r != nil {
		//Recovered from an error
		//todo Determine supervisor strategies and switch here
		//todo New callbacks in actor interface ?
		fmt.Println("panicing!!!")
	}
}

func (selfPtr *DefaultActor) handleError(err error) {

}

//This is the main runner of the actor
//The real stuff happens here
func (selfPtr *DefaultActor) run() {
	//Run recover in case an error occurs
	defer recoverFunc(selfPtr.actorInterface)
	defer selfPtr.selfLock.Unlock()

	if selfPtr.atomicallyCheckPtr(selfPtr.stopped) == 1 {
		err := selfPtr.actorInterface.OnDeadletter(convertDefaultActorToActorRef(selfPtr))
		if err != nil {
			panic(err)
		}
		return
	}

	//Lock self, so multiple threads can run this actor simultaneously (aka executors)
	selfPtr.selfLock.Lock()

	//Actor just started run the callback
	if selfPtr.atomicallyCheckPtr(selfPtr.justStarted) == 1 {
		err := selfPtr.actorInterface.OnStart(convertDefaultActorToActorRef(selfPtr))
		if err != nil {
			panic(err)
		}
		selfPtr.atomicallyResetPtr(selfPtr.justStarted)
	}

	//todo A better way to stop?
	if selfPtr.atomicallyCheckPtr(selfPtr.stop) == 1 {
		err := selfPtr.actorInterface.OnStop(convertDefaultActorToActorRef(selfPtr))
		if err != nil {
			panic(err)
		}
		selfPtr.atomicallySetPtr(selfPtr.stopped)
		return
	}

	//A pending stop was fired
	//There is only 1 item left in the message box, and a pending stop is active
	//So the message will be processed and stopped will be set
	if selfPtr.atomicallyCheckPtr(selfPtr.pendingStop) == 1 && selfPtr.getMessageQueueCount() == 1 {
		//set stopped to 1 if it is 0
		selfPtr.atomicallySetPtr(selfPtr.stopped)
	}

	item, check := selfPtr.messageQueue.GetOne()
	//There is no way check to be false, but still...
	if check {
		actorMessage := item.(ActorMessage)
		fmt.Println("Got message for actor : " + selfPtr.Name)

		fmt.Println("now total size of actor " + selfPtr.Name + "'s channel is reduced to " + strconv.Itoa(selfPtr.messageQueue.GetTotalItemCount()))
		err := selfPtr.actorInterface.OnReceive(convertDefaultActorToActorRef(selfPtr), actorMessage)
		if err != nil {
			panic(err)
		}
	}

}


/**
Starts the actor
 */
func (selfPtr *DefaultActor) Start() {
	//If stopped is 1 make it 0
	selfPtr.atomicallyResetPtr(selfPtr.stop)
	selfPtr.atomicallyResetPtr(selfPtr.stopped)
	selfPtr.atomicallyResetPtr(selfPtr.pendingStop)
}

/**
Messages told to actors are always of this type
The user creates a custom message struct or whatever and sets it to Msg

Users are responsible to retrieve actual Msg from ActorMessage
 */
type ActorMessage struct {
	Msg    interface{}
	Teller ActorRef
}

/**
In order to tell an actor a message, this should be used
 */
func (selfPtr *DefaultActor) Tell(msg interface{}, tellerRef ActorRef) error {
	if selfPtr.atomicallyCheckPtr(selfPtr.pendingStop) == 0 && selfPtr.atomicallyCheckPtr(selfPtr.stop) == 0 {
		//todo Do we need actor message?
		var actorMessage ActorMessage = ActorMessage{
			Msg: msg,
			Teller: tellerRef,
		}

		//Here the type of the actual item into the queue is ActorMessage
		//But we give the groupName manually as typename, so the round-robin will actually work
		typeName := getTypeNameFromType(reflect.TypeOf(msg))
		selfPtr.messageQueue.Enlist(typeName, actorMessage)

		//Enlist itself to be ran for this message
		//todo Would this create overhead somehow?
		//todo Fairness, starvation problems ?
		selfPtr.ActorSystem.actorChannel <- selfPtr

		return nil
	}

	return fmt.Errorf("This actor is pending for stop or is stopped, this message will not be processed")
}


/**
Stops the actor after the current message is processed
Does not delete the actor
 */
func (selfPtr *DefaultActor) StopRightAway() {
	selfPtr.startStopLock.Lock()
	defer selfPtr.startStopLock.Unlock()

	selfPtr.atomicallySetPtr(selfPtr.stop)
}

/**
Stops the actor after all the messages in the box are processed
 */
func (selfPtr *DefaultActor) PendingStop() {
	selfPtr.startStopLock.Lock()
	defer selfPtr.startStopLock.Unlock()

	selfPtr.atomicallySetPtr(selfPtr.pendingStop)
}

/**
Restarts the actor right away
 */
func (selfPtr *DefaultActor) Restart() {
	selfPtr.startStopLock.Lock()
	defer selfPtr.startStopLock.Unlock()

	//todo implement?
}