package actor
import (
	"strings"
	"fmt"
	"sync"
)

type ActorSystem struct {
	actorMap   map[string]Actor
	channelMap map[string]chan ActorMessage
	rootActor  *DefaultActor
	waitGroup  *sync.WaitGroup

}

/**
Initialize the pointers within actor system struct
 */
func (actorSystem *ActorSystem) InitSystem() {
	actorSystem.rootActor = new(DefaultActor)
	actorSystem.rootActor.Name = "root"
	actorSystem.rootActor.ChildrenMap = make(map[string]*DefaultActor)
	actorSystem.channelMap = make(map[string]chan ActorMessage)
	actorSystem.actorMap = make(map[string]Actor)
	actorSystem.waitGroup = new(sync.WaitGroup)
	//todo wht to do with actor interface in master ?
}

func (actorSystem *ActorSystem) Run() {
	actorSystem.waitGroup.Wait()
}

/**
Parse a path and return path elements as a slice
 */
func pathParser(name string) []string {
	path := strings.Split(name, "/") //path is a slice of names
	//An absolute name was given
	if string(name[0]) == string("/") {
		path = path[1:len(path)]
	}

	return path
}


/**
Get the actor corresponding to the path, starting from the given DefaultActor pointer
Runs recursively
//todo Implement a nonrecursive version for the future
 */
func getActorFromNodeRecursively(defaultActor *DefaultActor, currentIndexOfPathElement int, pathSlice []string) (*DefaultActor, error) {
	nodeForPathElement, isFound := defaultActor.ChildrenMap[pathSlice[currentIndexOfPathElement]]
	if isFound {
		//This is the last part of the uri path
		if currentIndexOfPathElement == len(pathSlice) - 1 {
			return nodeForPathElement, nil
		}

		return getActorFromNodeRecursively(nodeForPathElement, currentIndexOfPathElement + 1, pathSlice)
	} else {
		return nil, ActorNotFound{} //todo Write custom error
	}
}

/**
Remove the name of the actor to be created from the path to target for finding parent of the future actor
e.g. /actor/parentActor/futureActor => /actor/parentActor
 */
func getPathUntilParentAndNameOfFutureActor(path string) ([]string, string) {
	pathSlice := pathParser(path)
	return pathSlice[0:len(pathSlice) - 1], pathSlice[len(pathSlice) - 1]
}


/**
Get parent actor of the actor corresponding to the path
This method is used with an actor path which is not yet created
Once the parent of the actor to be created is found, the new actor
is created as a child under that actor
 */
func getParentRecursively(defaultActor *DefaultActor, path string) (*DefaultActor, string, error) {
	pathSliceUntilFutureParentName, nameOfFutureActor := getPathUntilParentAndNameOfFutureActor(path)

	if len(pathSliceUntilFutureParentName) == 0 {
		return defaultActor, nameOfFutureActor, nil
	}

	actor, err := getActorFromNodeRecursively(defaultActor, 0, pathSliceUntilFutureParentName)
	if err != nil {
		return nil, "", err
	}
	return actor, nameOfFutureActor, nil
}

/**
Create actor
 */
func (actorSystem *ActorSystem) CreateActor(actor Actor, path string) (ActorRef, error) {
	fmt.Println("create actor " + path)
	parentActor, singularName, err := getParentRecursively(actorSystem.rootActor, path)
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
	parentActor.ChildrenMap[singularName] = newActor

	//Create the listening channel for the new actor
	channelForActor := make(chan ActorMessage)

	//Create the stop channel for the actor
	stopChannelForActor := make(chan uint8)

	//Add the channel of the actor to the actor system channel map, with full path name
	actorSystem.channelMap[path] = channelForActor

	//Set the created channels
	newActor.Channel = channelForActor
	newActor.StopChannel = stopChannelForActor

	//Set actorSystems wait group to actor
	newActor.waitGroup = actorSystem.waitGroup

	//Increment wait group counter for newly created actor
	actorSystem.waitGroup.Add(1)

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
func (actorSystem *ActorSystem) GetActorRef(path string) (ActorRef, error) {
	rootActor := actorSystem.rootActor
	foundActor, err := getActorFromNodeRecursively(rootActor, 0, pathParser(path))
	if err != nil {
		return ActorRef{}, err
	}
	actorRefForFoundActor := convertDefaultActorToActorRef(foundActor)
	return actorRefForFoundActor, nil

}

