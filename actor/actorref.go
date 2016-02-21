package actor

var actorIndexer []*DefaultActor

type ActorRef struct {
	actorIndex int
}

func (actorRef *ActorRef) Tell(msg interface{}, tellerRef *ActorRef) {
	var defaultActor *DefaultActor = actorIndexer[actorRef.actorIndex]
	defaultActor.Tell(msg, tellerRef)
}
