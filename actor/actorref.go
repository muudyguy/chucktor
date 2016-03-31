package actor
import "reflect"

var actorIndexer []*CoreActor

//todo DID NOT DECIDE YET WHETHER TO MAKE PEOPLE USE ACTORREF AS POINTER OR VALUE
//todo RIGHT NOW CORRECT THING TO DO LOOKS LIKE MAKE THEM USE IT AS VALUE, CHANGE IT LATER?

type ActorRef struct {
	actorIndex int
	coreActor  *CoreActor
}

func (selfPtr *ActorRef) Tell(msg interface{}, tellerRef ActorRef) error {
//	var defaultActor *DefaultActor = actorIndexer[actorRef.actorIndex]
	return selfPtr.coreActor.Tell(msg, tellerRef)
}

func (selfPtr *ActorRef) Name() string {
	return selfPtr.coreActor.Name
}

func (selfPtr *ActorRef) Watch(requester ActorRef) {
	selfPtr.coreActor.Watch(requester.coreActor)
}

func (selfPtr *ActorRef) SetPriority(typ reflect.Type, priority int) {
	selfPtr.coreActor.SetMessagePriority(typ, priority)
}


//todo Maybe actorRefs should be created and saved into defaultactors in advance
//todo so we can avoid iterating and creating an actorref
func (actorRef *ActorRef) Children() []ActorRef {
	var childrenArray []*CoreActor = actorRef.coreActor.ChildrenArray
	var actorRefSlice []ActorRef
	for i := 0; i < len(childrenArray); i++ {
		actorRefSlice = append(actorRefSlice, convertCoreActorToActorRef(childrenArray[i]))
	}
	return actorRefSlice
}

func (selfPtr *ActorRef) Stop() {
	selfPtr.coreActor.StopNow()
}

func (selfPtr *ActorRef) ReStart() {
	selfPtr.coreActor.Restart()
}

/**
Create actor with path, using the current actor as parent
 */
func (selfPtr *ActorRef) CreateActor(actor ActorInterface, path string) {
	createActorOnParent(actor, selfPtr.coreActor.ActorSystem, path, selfPtr.coreActor)
}

/**
Get parent of this actor
 */
func (selfPtr *ActorRef) Parent() *ActorRef {
	return &ActorRef{
		coreActor:selfPtr.coreActor.Parent,
	}
}

//Tells to every children of the actor
//Returns the error ActorIsStopped or nil
//todo If a tell operation fails after the first child, some messages will be sent others wont
//todo But the error could only be the case when the child actor is stopped. Does it really matter ?
func (selfPtr *ActorRef) TellChildren(msg interface{}, tellerRef ActorRef) error {
	for _, childActor := range selfPtr.coreActor.ChildrenArray {
		actorIsStoppedError := childActor.Tell(msg, tellerRef)
		if actorIsStoppedError != nil {
			return actorIsStoppedError
		}
	}
	return nil
}

//Tells to the parent of the actor
//Returns the error ActorIsStopped or nil
func (selfPtr *ActorRef) TellParent(msg interface{}, tellerRef ActorRef) error {
	actorIsStoppedError := selfPtr.coreActor.Parent.Tell(msg, tellerRef)
	if actorIsStoppedError != nil {
		return actorIsStoppedError
	}
	return nil
}