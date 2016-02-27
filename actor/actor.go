package actor
import (

	"fmt"
//	"reflect"
	"sync/atomic"
	"reflect"
)

/**
A new actor must implement this to be created
 */
type Actor interface {
	OnReceive(self ActorRef, msg ActorMessage)
	OnStart(self ActorRef)
	OnStop(self ActorRef)
}


/**
Actual actors kept within the system
Childrens are kept in a slice, also a map
This brings memory overhead, but must be a bit faster if it was to be kept
in a map only
 */
type DefaultActor struct {
	Name           string
	ChildrenArray  []*DefaultActor
	ChildrenMap    map[string]*DefaultActor
	Parent         *DefaultActor
	actorInterface Actor
	Channel        PriorityBasedChannel
	index          int

	stopped        *uint32
	justStarted	   *uint32
}

func NewDefaultActor(name string, parent *DefaultActor) *DefaultActor {
	da := DefaultActor{}
	da.Name = name
	da.Parent = parent

	da.stopped = new(uint32)
	*da.stopped = 0

	da.justStarted = new(uint32)
	*da.justStarted = 1

	da.Channel = NewPriorityBasedChannel()

	da.Channel.setPriorityForMessageType(reflect.TypeOf(ActorMessage{}), 1)

	return &da
}

/**
This method stays alive as long as actor is alive
 */
func (selfPtr *DefaultActor) runner() {
	if atomic.LoadUint32(selfPtr.justStarted) == 1 {
		selfPtr.actorInterface.OnStart(convertDefaultActorToActorRef(selfPtr))
		atomic.CompareAndSwapUint32(selfPtr.justStarted, 1, 0)
	}

	fmt.Println("Starting message box for actor " + selfPtr.Name)
	for {
		actorMessage := selfPtr.Channel.Get()
		if atomic.LoadUint32(selfPtr.stopped) == 1 {
			selfPtr.actorInterface.OnStop(convertDefaultActorToActorRef(selfPtr))
			break
		}
		selfPtr.actorInterface.OnReceive(convertDefaultActorToActorRef(selfPtr), actorMessage.(ActorMessage))
	}

	fmt.Println("Stopped for : " + selfPtr.Name)
}

/**
Starts the actor worker
 */
func (defaultActor *DefaultActor) Start() {
	go defaultActor.runner()
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

/**
Stops the actor, but does not delete it
 */
func (selfPtr *DefaultActor) Stop() {
	if atomic.LoadUint32(selfPtr.stopped) == 0 {
		atomic.CompareAndSwapUint32(selfPtr.stopped, 0, 1)
	}
}