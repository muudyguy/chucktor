package actor
import (
	"muddle/tree"
	"fmt"
	"sync"
//	"reflect"
)

/**
A new actor must implement this to be created
 */
type Actor interface {
	OnReceive(self ActorRef, msg ActorMessage)
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
	Channel        chan ActorMessage
	StopChannel    chan uint8
	waitGroup      *sync.WaitGroup
	index          int
}

/**
This method stays alive as long as actor is alive
 */
func (defaultActor *DefaultActor) runner() {
	defer defaultActor.waitGroup.Done()

	var stop bool = false
	fmt.Println("Starting message box for actor " + defaultActor.Name)
	for {
		select {
		case <-defaultActor.StopChannel:
			stop = true
		default:
		}

		if !stop {
			actorMessage := <-defaultActor.Channel
			defaultActor.actorInterface.OnReceive(convertDefaultActorToActorRef(defaultActor), actorMessage)
		} else {
			fmt.Println("Stopping for : " + defaultActor.Name)
			break
		}

	}

	fmt.Println("Stopped for : " + defaultActor.Name)
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
	Msg    interface{}
	Teller ActorRef
}

/**
In order to tell an actor a message, this should be used
 */
func (defaultActor *DefaultActor) Tell(msg interface{}, tellerRef ActorRef) {
	var actorMessage ActorMessage = ActorMessage{
		Msg: msg,
		Teller: tellerRef,
	}
	defaultActor.Channel <- actorMessage
}

/**
Stops the actor, but does not delete it
//todo When do we delete it ?
 */
func (defaultActor *DefaultActor) Stop() {
	defaultActor.StopChannel <- 1
}

func (defaultActor DefaultActor) Bigger(daInterface tree.Comparable) bool {
	toBeCompared := daInterface.(DefaultActor)
	return defaultActor.Name > toBeCompared.Name
}

func (defaultActor DefaultActor) Equals(daInterface tree.Comparable) bool {
	toBeCompared := daInterface.(DefaultActor)
	return defaultActor.Name == toBeCompared.Name
}