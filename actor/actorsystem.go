package actor
import (
	"strings"
	"fmt"
)

type ActorSystem struct {
	actorMap map[string]Actor
	channelMap map[string]chan ActorMessage
	rootActor *DefaultActor

}

func (actorSystem *ActorSystem) InitSystem() {
	actorSystem.rootActor = new(DefaultActor)
	actorSystem.rootActor.Name = "root"
	actorSystem.rootActor.ChildrenMap = make(map[string]*DefaultActor)
	actorSystem.channelMap = make(map[string]chan ActorMessage)
	actorSystem.actorMap = make(map[string]Actor)
	//todo wht to do with actor interface in master ?
}


func nameParser(name string) []string {
	path := strings.Split(name, "/") //path is a slice of names
	//An absolute name was given
	if string(name[0]) == string("/") {
		path = path[1:len(path)]
	}

	return path
}

/**
Recursively searches until dst in filled with parent to be
But Even if one branch finds the actor, it keeps searching.

//todo In the future time complexity should be reduced by keeping a tree for names/actorPointers
 */
func recursivelyCheckForAvailabilityAndFillIfAvailable(actor *DefaultActor, currentIndexForNames int, pathSlice []string, dstParent **DefaultActor) {
	if actor.Name == pathSlice[currentIndexForNames] {
		if (currentIndexForNames == len(pathSlice) - 2) {
			*dstParent = actor
			return
		}
		for i := 0; i < len(actor.ChildrenArray); i++ {
			recursivelyCheckForAvailabilityAndFillIfAvailable(actor.ChildrenArray[i], currentIndexForNames + 1, pathSlice, dstParent);
		}
	} else {

	}
}

func (actorSystem *ActorSystem) getParent(name string) (string, *DefaultActor, error) {
	namesSlice := nameParser(name)
	rootActor := actorSystem.rootActor

	var dstParent **DefaultActor = new(*DefaultActor)
	if rootActor.Name == namesSlice[0] {
		return "", nil, fmt.Errorf("Cannot use the name of the internal root actor")
	} else {

		//If name is singular, parent is root
		if len(namesSlice) == 1 {
			*dstParent = rootActor
		} else {
			for i := 0; i < len(rootActor.ChildrenArray); i++ {
				recursivelyCheckForAvailabilityAndFillIfAvailable(rootActor.ChildrenArray[i], 0, namesSlice, dstParent)
			}

			if dstParent == nil {
				return "", nil, fmt.Errorf("Cannot find a viable parent")
			}
		}

		return namesSlice[len(namesSlice) - 1], *dstParent, nil
	}

}




func (actorSystem *ActorSystem) CreateActor(actor Actor, name string) (ActorRef, error) {
	singularName, parentActor, err := actorSystem.getParent(name)
	if err != nil {
		return ActorRef{}, err
	}

	//Create new actor
	var newActor *DefaultActor = new(DefaultActor)
	newActor.Name = singularName
	newActor.actorInterface = actor
	newActor.Parent = parentActor
	newActor.ChildrenMap = make(map[string]*DefaultActor)

	if parentActor == nil {
		return ActorRef{}, fmt.Errorf("In correct path, parent does not exist")
	}
	//Append the new actors pointer to parents children array
	appended := append(parentActor.ChildrenArray, newActor)
	parentActor.ChildrenArray = appended

	//Add new actor to the parent children map
	fmt.Println(parentActor.ChildrenMap)
	parentActor.ChildrenMap[singularName] = newActor

	//Create the listening channel for the new actor
	channelForActor := make(chan ActorMessage)

	//Create the stop channel for the actor
	stopChannelForActor := make(chan uint8)

	//Add the channel of the actor to the actor system channel map, with full path name
	actorSystem.channelMap[name] = channelForActor

	//Set the created channels
	newActor.Channel = channelForActor
	newActor.StopChannel = stopChannelForActor

	//Add the new actor pointer to the indexer for actor ref
	//todo Probably not needed anymore !
	actorIndexer = append(actorIndexer, newActor)

	newActor.index = len(actorIndexer) - 1

	newActor.Start()

	return ActorRef{
		actorIndex: newActor.index, //todo Not needed anymore ?
		defaultActor:newActor,
	}, nil
}

/**
Get actor from system with full path
 */
func (actorSystem *ActorSystem) GetActorRef(name string) (ActorRef , error) {
	singularName, parentActor, err := actorSystem.getParent(name)
	if err != nil {
		return ActorRef{}, err
	}
	actor := parentActor.ChildrenMap[singularName]
	return ActorRef{
		defaultActor:actor,
	}, nil
}

