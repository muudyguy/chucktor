package actor
import "fmt"

/**
A new actor must implement this to be created
 */
type ActorInterface interface {
	OnRecieve(self *ActorRef, msg ActorMessage)
}

/**
Actual actors kept within the system
 */
type DefaultActor struct {
	Name string
	Children []*DefaultActor
	Parent *DefaultActor
	actorInterface ActorInterface
	Channel chan ActorMessage
	StopChannel chan uint8
	index int
}

/**
This method stays alive as long as actor is alive
 */
func (defaultActor *DefaultActor) runner() {
	fmt.Println("started")
	var stop bool = false
	for {
		select {
		case <- defaultActor.StopChannel:
			stop = true
		default:
		}

		if !stop {
//			fmt.Println("waiting message")
			actorMessage := <- defaultActor.Channel
//			fmt.Println("got message")
			defaultActor.actorInterface.OnRecieve(convertDefaultActorToActorRef(defaultActor), actorMessage)
		} else {
			break
		}

	}
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
	Teller *ActorRef
}

/**
In order to tell an actor a message, this should be used
 */
func (defaultActor *DefaultActor) Tell(msg interface{}, tellerRef *ActorRef) {
	var actorMessage ActorMessage = ActorMessage{
		Msg: msg,
		Teller: tellerRef,
	}
	defaultActor.Channel <- actorMessage
}


func (defaultActor *DefaultActor) Stop() {
	defaultActor.StopChannel <- 1
}