package actor

func convertDefaultActorToActorRef(defaultActor *DefaultActor) *ActorRef {
	actorRef := ActorRef{defaultActor.index, defaultActor}
	return &actorRef
}
