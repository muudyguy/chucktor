package actor
import "reflect"

var actorIndexer []*CoreActor

//todo DID NOT DECIDE YET WHETHER TO MAKE PEOPLE USE ACTORREF AS POINTER OR VALUE
//todo RIGHT NOW CORRECT THING TO DO LOOKS LIKE MAKE THEM USE IT AS VALUE, CHANGE IT LATER?

type ActorRef struct {
	actorIndex int
	defaultActor *CoreActor
}

func (selfPtr *ActorRef) Tell(msg interface{}, tellerRef ActorRef) error {
//	var defaultActor *DefaultActor = actorIndexer[actorRef.actorIndex]
	return selfPtr.defaultActor.Tell(msg, tellerRef)
}

func (selfPtr *ActorRef) SetPriority(typ reflect.Type, priority int) {
	selfPtr.defaultActor.SetMessagePriority(typ, priority)
}


//todo Maybe actorRefs should be created and saved into defaultactors in advance
//todo so we can avoid iterating and creating an actorref
func (actorRef *ActorRef) Children() []ActorRef {
	var childrenArray []*CoreActor = actorRef.defaultActor.ChildrenArray
	var actorRefSlice []ActorRef
	for i := 0; i < len(childrenArray); i++ {
		actorRefSlice = append(actorRefSlice, convertDefaultActorToActorRef(childrenArray[i]))
	}
	return actorRefSlice
}

func (selfPtr *ActorRef) Stop() {
	selfPtr.defaultActor.StopRightAway()
}

func (selfPtr *ActorRef) ReStart() {
	selfPtr.defaultActor.Restart()
}

/**
Create actor with path, using the current actor as parent
 */
func (selfPtr *ActorRef) CreateActor(actor ActorInterface, path string) {
	createActorOnParent(actor, selfPtr.defaultActor.ActorSystem, path, selfPtr.defaultActor)
}

/**
Get parent of this actor
 */
func (selfPtr *ActorRef) Parent() *ActorRef {
	return &ActorRef{
		defaultActor:selfPtr.defaultActor.Parent,
	}
}