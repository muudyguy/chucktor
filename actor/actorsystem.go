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
	actorSystem.channelMap = make(map[string]chan ActorMessage)
	actorSystem.actorMap = make(map[string]Actor)
	//todo wht to do with actor interface in master ?
}


func nameParser(name string) []string {
	path := strings.Split(name, "/") //path is a slice of names
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
		for i := 0; i < len(actor.Children); i++ {
			recursivelyCheckForAvailabilityAndFillIfAvailable(actor.Children[i], currentIndexForNames + 1, pathSlice, dstParent);
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
			for i := 0; i < len(rootActor.Children); i++ {
				recursivelyCheckForAvailabilityAndFillIfAvailable(rootActor.Children[i], 0, namesSlice, dstParent)
			}

			if dstParent == nil {
				return "", nil, fmt.Errorf("Cannot find a viable parent")
			}
		}

		return namesSlice[len(namesSlice) - 1], *dstParent, nil
	}

}




func (actorSystem *ActorSystem) CreateActor(actor Actor, name string) *ActorRef {
	singularName, parentActor, err := actorSystem.getParent(name)
	if err != nil {
		fmt.Println(err)
	}

	newActor := new(DefaultActor)
	newActor.Name = singularName
	newActor.actorInterface = actor
	newActor.Parent = parentActor
	appended := append(parentActor.Children, newActor)
	parentActor.Children = appended

	channelForActor := make(chan ActorMessage)
	stopChannelForActor := make(chan uint8)
	actorSystem.channelMap[name] = channelForActor
	newActor.Channel = channelForActor
	newActor.StopChannel = stopChannelForActor

	actorIndexer = append(actorIndexer, newActor)

	newActor.index = len(actorIndexer) - 1

	newActor.Start()

	return &ActorRef{
		actorIndex: newActor.index,
	}
}

