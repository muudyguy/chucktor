package actor

var actorIndexer []*DefaultActor

type ActorRef struct {
	actorIndex int
	defaultActor *DefaultActor
}

func (actorRef *ActorRef) Tell(msg interface{}, tellerRef *ActorRef) {

//	var defaultActor *DefaultActor = actorIndexer[actorRef.actorIndex]
	actorRef.defaultActor.Tell(msg, tellerRef)
}
