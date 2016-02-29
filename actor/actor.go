package actor
import (

	"fmt"
	"sync/atomic"
	"reflect"
	"strconv"
	"sync"
)

/**
A new actor must implement this to be created
 */
type Actor interface {
	OnReceive(self ActorRef, msg ActorMessage)
	OnStart(self ActorRef)
	OnStop(self ActorRef)
	OnRestart(self ActorRef)
}


/**
Actual actors kept within the system
Children are kept in a slice, also a map
This brings memory overhead, but must be a bit faster if it was to be kept
in a map only
 */
type DefaultActor struct {
	Name             string
	ChildrenArray    []*DefaultActor
	ChildrenMap      map[string]*DefaultActor
	Parent           *DefaultActor
	actorInterface   Actor
	Channel          PriorityBasedChannel
	index            int

	stopped          *uint32

	restarted        *uint32
	restartedChannel chan int
	stoppedChannel   chan int
	justStarted      *uint32

	startStopLock    *sync.Mutex
}

func NewDefaultActor(name string, parent *DefaultActor) *DefaultActor {
	da := DefaultActor{}
	da.Name = name
	da.Parent = parent

	da.stopped = new(uint32)
	*da.stopped = 0

	da.justStarted = new(uint32)
	*da.justStarted = 1

	da.restarted = new(uint32)
	*da.restarted = 0

	da.Channel = NewPriorityBasedChannel(name + " channel")

	da.Channel.setPriorityForMessageType(reflect.TypeOf(ActorMessage{}), 1)

	da.restartedChannel = make(chan int)
	da.stoppedChannel = make(chan int, 1)

	da.startStopLock = new(sync.Mutex)

	return &da
}

/**
This method stays alive as long as actor is alive
 */
func (selfPtr *DefaultActor) runner(restarted bool) {
	if restarted {
		selfPtr.actorInterface.OnRestart(convertDefaultActorToActorRef(selfPtr))
	} else {
		if atomic.LoadUint32(selfPtr.justStarted) == 1 {
			selfPtr.actorInterface.OnStart(convertDefaultActorToActorRef(selfPtr))
			atomic.CompareAndSwapUint32(selfPtr.justStarted, 1, 0)
		}
	}

	for {

		//todo A better way to stop?
		if atomic.LoadUint32(selfPtr.stopped) == 1 {
			selfPtr.actorInterface.OnStop(convertDefaultActorToActorRef(selfPtr))
			selfPtr.stoppedChannel <- 1
			break
		}

		actorMessage := selfPtr.Channel.Get()

		switch actorMessage.(type) {
		case DummyMessage:
			continue
		//StoppedAfterQueueComplete
		case StopMessage:
			selfPtr.actorInterface.OnStop(convertDefaultActorToActorRef(selfPtr))
			break
		}
		fmt.Println("Got message for actor : " + selfPtr.Name)
		fmt.Println(actorMessage)
		fmt.Println("now total size of actor " + selfPtr.Name + "'s channel is reduced to " + strconv.Itoa(selfPtr.Channel.messageQueue.GetTotalItemCount()))
		selfPtr.actorInterface.OnReceive(convertDefaultActorToActorRef(selfPtr), actorMessage.(ActorMessage))
	}


}

/**
Starts the actor worker
 */
func (selfPtr *DefaultActor) Start() {
	//If stopped is 1 make it 0
	atomic.CompareAndSwapUint32(selfPtr.stopped, 1, 0)
	go selfPtr.runner(false)
}

func (selfPtr *DefaultActor) restart() {
	//If stopped is 1 make it 0
	atomic.CompareAndSwapUint32(selfPtr.stopped, 1, 0)
	go selfPtr.runner(true)
}


/**
Messages told to actors are always of this type
The user creates a custom message struct or whatever and sets it to Msg

Users are responsible to retrieve actual Msg from ActorMessage
 */
type ActorMessage struct {
	Stop   bool
	Msg    interface{}
	Teller ActorRef
}

/**
In order to tell an actor a message, this should be used
 */
func (selfPtr *DefaultActor) Tell(msg interface{}, tellerRef ActorRef) {
	var actorMessage ActorMessage = ActorMessage{
		Stop:false,
		Msg: msg,
		Teller: tellerRef,
	}

	selfPtr.Channel.Send(actorMessage)
}

type StopMessage struct  {

}


type DummyMessage struct {

}

/**
Stops the actor after the existing messages are processed
Does not delete the actor
 */
func (selfPtr *DefaultActor) StopAfterQueueComplete() {
	selfPtr.Tell(StopMessage{}, ActorRef{})
}

/**
Stops the actor after the current message is processed
Does not delete the actor
 */
func (selfPtr *DefaultActor) StopRightAway() {
	selfPtr.startStopLock.Lock()
	defer selfPtr.startStopLock.Unlock()
	fmt.Println("stopping")
	if atomic.LoadUint32(selfPtr.stopped) == 0 {
		fmt.Println("stop was 0")
		atomic.CompareAndSwapUint32(selfPtr.stopped, 0, 1)
		selfPtr.Tell(DummyMessage{}, ActorRef{}) //If there are no messages send a dummy message???

		//This makes stop stall, maybe run this with goRoutine ?

		<- selfPtr.stoppedChannel

	}
	fmt.Println("Stop released the lock")
}

/**
Restarts the actor
 */
func (selfPtr *DefaultActor) Restart() {
	selfPtr.startStopLock.Lock()
	defer selfPtr.startStopLock.Unlock()

	fmt.Println("Started restarting!")

	if atomic.LoadUint32(selfPtr.stopped) == 0 {
		atomic.CompareAndSwapUint32(selfPtr.stopped, 0, 1)
		selfPtr.Tell(nil, ActorRef{}) //If there are no messages send a dummy message???

		<- selfPtr.stoppedChannel
		fmt.Println("stopped")

	}

	selfPtr.restart()

	fmt.Println("restarted")
}