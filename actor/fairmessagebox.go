package actor
import (
	"sync"
	"fmt"
)

type ActorState struct {
	currentQuantum int
	messageCount int
	coreActor *CoreActor
}

func NewActorState(coreActor *CoreActor) *ActorState {
	as := ActorState{}
	as.coreActor = coreActor
	return &as
}

//Has to be used under strict locking
type RollableActorStateSlice []*ActorState

func (selfPtr *RollableActorStateSlice) GetFront() *ActorState {
	return []*ActorState(*selfPtr)[0]
}

//todo Does it copy the content ?
func (selfPtr *RollableActorStateSlice) Roll() {
	sliceSelf := []*ActorState(*selfPtr)
	length := len(sliceSelf)
	*selfPtr = append(sliceSelf[1:length], sliceSelf[0])
}

func (selfPtr *RollableActorStateSlice) Remove(actorState *ActorState) {
	sliceSelf := []*ActorState(*selfPtr)
	length := len(sliceSelf)
	for i, value := range sliceSelf {
		if value == actorState {
			*selfPtr = append(sliceSelf[0:i], sliceSelf[i + 1:length]...)
		}
	}
}

func (selfPtr *RollableActorStateSlice) Push(actorState *ActorState) {
	sliceSelf := []*ActorState(*selfPtr)
	*selfPtr = append(sliceSelf, actorState)
}

type FairGlobalMessageBox struct {
	//The fact that key is the full path of the actor might make this quite unviable if the user has very deeply nested
	//Actor hierarchy, a solution is to create unique ids for actors, but with which algorithm?
	actorStatesMap map[string]*ActorState
	quantum int
	actorStateQueue RollableActorStateSlice

	messageBoxLock *sync.Mutex
}

func NewFairGlobalMessageBox(quantum int) *FairGlobalMessageBox {
	fgmb := FairGlobalMessageBox{}

	fgmb.quantum = quantum
	fgmb.actorStateQueue = make([]*ActorState, 0, 100) //Make initial 100 capacity. Is it good though ?
	fgmb.messageBoxLock = new(sync.Mutex)
	fgmb.actorStatesMap = make(map[string]*ActorState)

	return &fgmb
}


func (selfPtr *FairGlobalMessageBox) Pop() *CoreActor {
	defer selfPtr.messageBoxLock.Unlock()
	selfPtr.messageBoxLock.Lock()

	actorState := selfPtr.actorStateQueue.GetFront()
	actorState.currentQuantum += 1
	actorState.messageCount -= 1

	if actorState.currentQuantum == selfPtr.quantum {
		actorState.currentQuantum = 0
		selfPtr.actorStateQueue.Roll()
	}

	if actorState.messageCount == 0 {
		selfPtr.actorStateQueue.Remove(actorState)

		//todo Is deleting absolutely necessary ?
		delete(selfPtr.actorStatesMap, actorState.coreActor.FullPath)
	}

	return actorState.coreActor
}

func (selfPtr *FairGlobalMessageBox) Enlist(coreActor *CoreActor) {
	defer selfPtr.messageBoxLock.Unlock()
	selfPtr.messageBoxLock.Lock()


	value, ok := selfPtr.actorStatesMap[coreActor.FullPath]
	if !ok {
		as := NewActorState(coreActor)
		fmt.Println(as)
		as.messageCount = 1
		selfPtr.actorStatesMap[coreActor.FullPath] = as
		selfPtr.actorStateQueue.Push(as)
	} else {
		value.messageCount += 1
	}
}