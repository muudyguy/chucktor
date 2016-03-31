package actor
import (
	"strings"
	"fmt"
	"sync"
	"queue"
)

type ActorSystem struct {
	actorMap     map[string]ActorInterface
	rootActor    *CoreActor
	actorChannel chan int

	executorNumber int
	executors []*ActorSystemExecutor

	messageBox MessageBoxInterface
}

func NewActorSystem(executorNumber int) *ActorSystem {
	as := &ActorSystem{}
	as.rootActor = newRootActor(as)
	as.actorMap = make(map[string]ActorInterface)
	as.createAndStartExecutors(executorNumber)
	as.actorChannel = make(chan int, 1000)
	as.messageBox = NewFairGlobalMessageBox(10)
	return as
}

/**
Initialize the pointers within actor system struct
 */
func (selfPtr *ActorSystem) InitSystem() {
	selfPtr.rootActor = newRootActor(selfPtr)
	selfPtr.actorMap = make(map[string]ActorInterface)
}

func newRootActor(actorSystem *ActorSystem) *CoreActor {
	da := CoreActor{}
	da.Name = "root"
	da.Parent = nil

	da.stop = new(uint32)
	*da.stop = 0

	da.justStarted = new(uint32)
	*da.justStarted = 1

	da.ActorSystem = actorSystem

	da.ChildrenArray = make([]*CoreActor, 0, 20) // How much default cap ?

	da.ChildrenMap = make(map[string]*CoreActor)

	da.messageQueue = queue.NewRoundRobinQueue()

	da.actorInterface = nil //todo Implement a default root interface for messages

	da.startStopLock = new(sync.Mutex)

	return &da
}

//Create and start the executors responsible to process actor messages
func (selfPtr *ActorSystem) createAndStartExecutors(executorNumber int) {
	for i := 0; i < executorNumber; i++ {
		executor := newActorSystemExecutor(selfPtr)
		selfPtr.executors = append(selfPtr.executors, executor)
		go executor.startExecution()
	}
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
func getActorFromNodeRecursively(defaultActor *CoreActor, currentIndexOfPathElement int,
		pathSlice []string) (*CoreActor, error) {

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
func getParentRecursively(defaultActor *CoreActor, path string) (*CoreActor, string, error) {
	pathSliceUntilFutureParentName, nameOfFutureActor := getPathUntilParentAndNameOfFutureActor(path)
//	fmt.Println("NAMES")
//	fmt.Println(path)
//	fmt.Println(pathSliceUntilFutureParentName)
//	fmt.Println(nameOfFutureActor)

	if len(pathSliceUntilFutureParentName) == 0 {
		return defaultActor, nameOfFutureActor, nil
	}

	actor, err := getActorFromNodeRecursively(defaultActor, 0, pathSliceUntilFutureParentName)
//	fmt.Println(actor.actorInterface)
	if err != nil {
		return nil, "", err
	}
	return actor, nameOfFutureActor, nil
}

/**

 */
func createActorOnParent(actor ActorInterface, actorSystem *ActorSystem, path string,
		parentStartActor *CoreActor) (ActorRef, error) {
	//Create new actor
	parentActor, singularName, err := getParentRecursively(parentStartActor, path)

	if err != nil {
		return ActorRef{}, err
	}

	var newActor *CoreActor = NewCoreActor(singularName, parentActor)
	newActor.FullPath = path //todo Move this into new core actor

	newActor.actorInterface = actor

	newActor.ChildrenMap = make(map[string]*CoreActor)

	if parentActor == nil {
		return ActorRef{}, fmt.Errorf("In correct path, parent does not exist")
	}
	//Append the new actors pointer to parents children array
	appended := append(parentActor.ChildrenArray, newActor)
	parentActor.ChildrenArray = appended

	//Add new actor to the parent children map
	parentActor.ChildrenMap[singularName] = newActor

	//Add the new actor pointer to the indexer for actor ref
	//todo Probably not needed anymore !
	actorIndexer = append(actorIndexer, newActor)

	newActor.index = len(actorIndexer) - 1

	return ActorRef{
		actorIndex: newActor.index, //todo Not needed anymore ?
		coreActor:newActor,
	}, nil
}

/**
Create actor
 */
func (selfPtr *ActorSystem) CreateActor(actor ActorInterface, path string) (ActorRef, error) {
	return createActorOnParent(actor, selfPtr, path, selfPtr.rootActor)
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
	actorRefForFoundActor := convertCoreActorToActorRef(foundActor)
	return actorRefForFoundActor, nil

}

