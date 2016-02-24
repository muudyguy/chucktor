package actor

var actorIndexer []*DefaultActor

type ActorRef struct {
	actorIndex int
	defaultActor *DefaultActor
}

func (actorRef *ActorRef) Tell(msg interface{}, tellerRef ActorRef) {

//	var defaultActor *DefaultActor = actorIndexer[actorRef.actorIndex]
	actorRef.defaultActor.Tell(msg, tellerRef)
}


//todo Maybe actorRefs should be created and saved into defaultactors in advance
//todo so we can avoid iterating and creating an actorref
func (actorRef *ActorRef) Children() []ActorRef {
	var childrenArray []*DefaultActor = actorRef.defaultActor.ChildrenArray
	var actorRefSlice []ActorRef
	for i := 0; i < len(childrenArray); i++ {
		actorRefSlice = append(actorRefSlice, convertDefaultActorToActorRef(childrenArray[i]))
	}
	return actorRefSlice
}