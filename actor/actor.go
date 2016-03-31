package actor
import (

	"sync/atomic"
	"reflect"
//	"strconv"
	"sync"
	"queue"
//	"strconv"
)



//todo FILE TOO BIG?
/**
Actual actors kept within the system
Children are kept in a slice, also a map
This brings memory overhead, but must be a bit faster if it was to be kept
in a map only
 */
type CoreActor struct {
	Name           string
	FullPath	   string
	ChildrenArray  []*CoreActor
	ChildrenMap    map[string]*CoreActor
	Watchers	   []*CoreActor
	Parent         *CoreActor
	actorInterface ActorInterface
	messageQueue   *queue.RoundRobinQueue
	index          int     //todo is it going to be used
	ActorSystem    *ActorSystem

	stop           *uint32 //todo Find atomic operations for 8 bit or boolean
	pendingStop    *uint32
	stoppingChannel chan int

	justStarted    *uint32

	restart        *uint32

	stopped        *uint32

	startStopLock  *sync.Mutex
	selfLock       *sync.Mutex

	supervisionStrategy int
}

/**
Create a new core actor with correct initializations
 */
func NewCoreActor(name string, parent *CoreActor) *CoreActor {
	da := CoreActor{}
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

	da.stoppingChannel = make(chan int, 10)

	da.ActorSystem = parent.ActorSystem

	da.messageQueue = queue.NewRoundRobinQueue()

	da.startStopLock = new(sync.Mutex)
	da.selfLock = new(sync.Mutex)

	da.supervisionStrategy = 0 //Propagate

	return &da
}

//Set priority for message
func (selfPtr *CoreActor) SetMessagePriority(typ reflect.Type, priority int) {
	selfPtr.messageQueue.SetGroup(getTypeNameFromType(typ), priority)
}

/**
This Block sets, resets and checks states
uint32 type could be changed for something else later using atomic.Value
I am not sure which one is more efficient. A benchmark could be good
 */

func (selfPtr *CoreActor) getMessageQueueCount() int {
	return selfPtr.messageQueue.GetTotalItemCount()
}

func (selfPtr *CoreActor) atomicallySetPtr(ptr *uint32) {
	atomic.CompareAndSwapUint32(ptr, 0, 1)
}

func (selftPtr *CoreActor) atomicallyResetPtr(ptr *uint32) {
	atomic.CompareAndSwapUint32(ptr, 1, 0)
}

func (selfPtr *CoreActor) atomicallyCheckPtr(ptr *uint32) uint32 {
	return atomic.LoadUint32(ptr)
}

func (selfPtr *CoreActor) serveSelf() {
	selfPtr.ActorSystem.messageBox.Enlist(selfPtr)
	selfPtr.ActorSystem.actorChannel <- 1
}


/**
Actor instance here will be a pointer as long as the user passes the pointer of the interface they implemented
//todo How to make them do it all the time to avoid mistakes in their behalf
 */
func (selfPtr *CoreActor) recoverFunc(actor ActorInterface) {
	if r := recover(); r != nil {
		//might fail if some other type was returned from recover()
		err := r.(error)
		selfPtr.handleError(err)
	}
}

/**
Stop all children of this core actor
 */
func (selfPtr *CoreActor) stopChildren() {
	for _, child := range selfPtr.ChildrenArray {
		child.StopNow()
	}
	//todo Delete, memory ?
	//how about selfPtr.Children = []*CoreActor ??  No, just 1 ActorRef would cause a memory leak. Maybe it should ?
}

/**
Handle any errors that occur in run()
 */
func (selfPtr *CoreActor) handleError(err error) {
	//todo This does not look good to me
	if selfPtr.Name == "root" {
		//todo do something else for root
		return
	}

	selfPtr.actorInterface.OnError(convertCoreActorToActorRef(selfPtr))
	switch selfPtr.supervisionStrategy {
	case 0:
		//todo Will this be enough to propagate
		selfPtr.Parent.handleError(err)
	case 1:
		//todo RESTART
	}
}

//This is the main runner of the actor
//The real stuff happens here
func (selfPtr *CoreActor) run() {
	//Run recover in case an error occurs
	//Right now executors receive actor state from the queue and run this method.
	//Only one actor state exists for an actor, so thre is no way for a simultaneous running of this method
	defer selfPtr.recoverFunc(selfPtr.actorInterface)
	defer selfPtr.selfLock.Unlock()
	//Lock self, so multiple threads can run this actor simultaneously (aka executors)
	//todo Do we really need to lock ?
	selfPtr.selfLock.Lock()


	if selfPtr.atomicallyCheckPtr(selfPtr.stopped) == 1 {
		err := selfPtr.actorInterface.OnDeadletter(convertCoreActorToActorRef(selfPtr))
		if err != nil {
			panic(err)
		}
		return
	}

	//Actor just started run the callback
	if selfPtr.atomicallyCheckPtr(selfPtr.justStarted) == 1 {
		err := selfPtr.actorInterface.OnStart(convertCoreActorToActorRef(selfPtr))
		if err != nil {
			panic(err)
		}
		selfPtr.atomicallyResetPtr(selfPtr.justStarted)
	}

	item, check := selfPtr.messageQueue.GetOne()


//	fmt.Println("Got message for actor : " + selfPtr.Name)
//	fmt.Println(item)
//	fmt.Println("now total size of actor " + selfPtr.Name + "'s channel is reduced to " + strconv.Itoa(selfPtr.messageQueue.GetTotalItemCount()))



	var err error
	if check {
		actorMessage := item.(ActorMessage)
		switch actorMessage.Msg.(type) {
			case StopNow:
				err = selfPtr.actorInterface.OnStop(convertCoreActorToActorRef(selfPtr))
				selfPtr.atomicallySetPtr(selfPtr.stopped)
				selfPtr.notifyWatchers() //Send necessary Death messages
//				fmt.Println("Received STOPNOW IN ACTOR :" + selfPtr.FullPath)
			case StopWhenDone:
				selfPtr.atomicallySetPtr(selfPtr.pendingStop)
				//todo Implement stop when done
				//todo The trick in my mind is to start a counter to process already existing messages and then set stopped flag
//				fmt.Println("Received STOPWHENDONE IN ACTOR :" + selfPtr.FullPath)
			case Restart:
				selfPtr.actorInterface.OnRestart(convertCoreActorToActorRef(selfPtr))
				//todo Notify message box to clear ? Or should the existing message processed?
				//todo Decide !
//				fmt.Println("Received RESTART IN ACTOR :" + selfPtr.FullPath)
			default:
				err = selfPtr.actorInterface.OnReceive(convertCoreActorToActorRef(selfPtr), actorMessage)
//				fmt.Println("Received DEFAULT IN ACTOR :" + selfPtr.FullPath)
		}

		if err != nil {
			panic(err)
		}
	} else {
		//todo do something in case
	}

}


/**
Starts the actor
 */
func (selfPtr *CoreActor) Start() {
	//If stopped is 1 make it 0
	selfPtr.atomicallyResetPtr(selfPtr.stopped)
	selfPtr.atomicallyResetPtr(selfPtr.pendingStop)

	selfPtr.serveSelf()
}



/**
In order to tell an actor a message, this should be used
 */
func (selfPtr *CoreActor) Tell(msg interface{}, tellerRef ActorRef) error {

	//todo When the actor is stopped, do not tell or process as dead letters?
	if selfPtr.atomicallyCheckPtr(selfPtr.pendingStop) == 0 && selfPtr.atomicallyCheckPtr(selfPtr.stopped) == 0 {
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
		//todo Too many hanging goroutines ?
		//todo Maybe instead of goroutines give huge buffer? It could work
		/**
		Another queue for executors could be designed. Priorities would be equal for every actor created.
		It could as well just be a prioritizing map, holding quantum states for every actor path
		if a quantum is zero, reset it and go to the next actor. Every Tell, also sends a message nevertheless.
		The problem of blocking could be solved with
		 */
		selfPtr.serveSelf()
		return nil
	}


	return ActorIsStopped{}
}


/**
Stops the actor after the current message is processed
Does not delete the actor
 */
func (selfPtr *CoreActor) StopNow() {
	selfPtr.startStopLock.Lock()
	defer selfPtr.startStopLock.Unlock()

	//todo What to do with the teller? We should enter a stopper
	selfPtr.messageQueue.EnlistAbsolutePriority(ActorMessage{Msg:StopNow{}})
	//Send itself to be ran for stop op
	selfPtr.serveSelf()
}

/**
Stops the actor after all the messages in the box are processed
 */
func (selfPtr *CoreActor) StopWhenDone() {
	selfPtr.startStopLock.Lock()
	defer selfPtr.startStopLock.Unlock()

	//todo What to do with the teller? We should enter a stopper
	selfPtr.messageQueue.EnlistAbsolutePriority(ActorMessage{Msg:StopWhenDone{}})
	selfPtr.serveSelf()
}

/**
Restarts the actor right away
 */
func (selfPtr *CoreActor) Restart() {
	selfPtr.startStopLock.Lock()
	defer selfPtr.startStopLock.Unlock()

	selfPtr.messageQueue.EnlistAbsolutePriority(ActorMessage{Msg:Restart{}})
	//Send itself to be ran for stop op
	selfPtr.serveSelf()
	//todo implement?
}

func (selfPtr *CoreActor) notifyWatchers() {
	for _, watcher := range selfPtr.Watchers {
		//todo This feels awkward here
		//todo It should be telling himself with a pointer ?
		watcher.Tell(Death{Dead:convertCoreActorToActorRef(selfPtr)}, convertCoreActorToActorRef(selfPtr))
	}
}

func (selfPtr *CoreActor) Watch(requester *CoreActor) {
	selfPtr.Watchers = append(selfPtr.Watchers, requester)
}