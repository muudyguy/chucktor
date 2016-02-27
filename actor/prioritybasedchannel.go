package actor
import (
	"reflect"
	"sync"
	"queue"
	"fmt"

	"sync/atomic"
)

type ChannelInterface interface {
	Send(msg interface{})
	Get() interface{}
}

type PriorityBasedChannel struct {
	messageArrivalChannel chan int

	priorityMap           *map[string]int
	messageQueue          *queue.RoundRobinQueue

	waiterCount           *uint32

	priorityMapLock       *sync.Mutex
}

func NewPriorityBasedChannel() PriorityBasedChannel {
	pbc := PriorityBasedChannel{}
	pbc.priorityMap = new(map[string]int)
	*pbc.priorityMap = make(map[string]int)

	pbc.priorityMapLock = new(sync.Mutex)

	pbc.messageQueue = queue.NewRoundRobinQueue()

	pbc.messageArrivalChannel = make(chan int)

	pbc.waiterCount = new(uint32)

	return pbc
}




func (selfPtr *PriorityBasedChannel) SetPriority(priority int, msgType reflect.Type) {
	selfPtr.priorityMapLock.Lock()
	defer selfPtr.priorityMapLock.Unlock()
	typeName := getTypeNameFromType(msgType)
	(*selfPtr.priorityMap)[typeName] = priority
	selfPtr.messageQueue.SetGroup(typeName, priority)
}

/**
Looks a bit messy and low performance
 */
func (selfPtr *PriorityBasedChannel) Send(msg interface{}) {

	selfPtr.priorityMapLock.Lock()
	if len(*selfPtr.priorityMap) == 0 {
		panic(fmt.Errorf("There are no priorities set within the map !"))
	}
	selfPtr.priorityMapLock.Unlock()

	groupName := getGroupNameFromMsg(msg)
	selfPtr.messageQueue.Enlist(groupName, msg)


	if atomic.LoadUint32(selfPtr.waiterCount) > 0 {
		selfPtr.messageArrivalChannel <- 1
	}



}

func (selfPtr *PriorityBasedChannel) availableMessage() interface{} {
	//queue is threadsafe
	//don't need the bool value
	item, _ := selfPtr.messageQueue.GetOne()
	return item
}

//Does not seem like it needs mutex locking
func (selfPtr *PriorityBasedChannel) Get() interface{} {
	selfPtr.priorityMapLock.Lock()
	if len(*selfPtr.priorityMap) == 0 {
		panic(fmt.Errorf("There are no priorities set within the map !"))
	}
	selfPtr.priorityMapLock.Unlock()

	message := selfPtr.availableMessage()
	if message == nil {
		atomic.AddUint32(selfPtr.waiterCount, 1)

		<-selfPtr.messageArrivalChannel

		atomic.AddUint32(selfPtr.waiterCount, ^uint32(0))

		message = selfPtr.availableMessage()

		return message
	}
	return message
}

func (selfPtr *PriorityBasedChannel) setPriorityForMessageType(typ reflect.Type, priority int) {
	typeName := getTypeNameFromType(typ)
	selfPtr.messageQueue.SetGroup(typeName, priority)

	selfPtr.priorityMapLock.Lock()
	defer selfPtr.priorityMapLock.Unlock()
	(*selfPtr.priorityMap)[typeName] = priority
}
