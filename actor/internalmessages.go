package actor

type Death struct {
	Dead ActorRef
}


type StopNow struct {

}

type StopWhenDone struct {

}

type Restart struct {

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

