package actor
import (
	"reflect"
	"sync"
	"queue"
	"fmt"

)

type ChannelInterface interface {
	Send(msg interface{})
	Get() interface{}
}

type PriorityBasedChannel struct {
	messageArrivalChannel chan int

	priorityMap           *map[string]int
	messageQueue          *queue.RoundRobinQueue

	waiterCount           *int

	waiterCountLock       sync.Mutex
	sendLock              sync.Mutex
	priorityMapLock       sync.Mutex
}

func (selfPtr *PriorityBasedChannel) Init() {
	selfPtr.priorityMap = new(map[string]int)
	*selfPtr.priorityMap = make(map[string]int)

	selfPtr.messageQueue = new(queue.RoundRobinQueue)

	selfPtr.messageArrivalChannel = make(chan int)

	selfPtr.waiterCount = new(int)
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
	//Avoid simultaneous/concurrent access of maps
	selfPtr.sendLock.Lock()
	defer selfPtr.sendLock.Unlock()

	selfPtr.priorityMapLock.Lock()
	if len(*selfPtr.priorityMap) == 0 {
		panic(fmt.Errorf("There are no priorities set within the map !"))
	}
	selfPtr.priorityMapLock.Unlock()

	groupName := getGroupNameFromMsg(msg)
	//	fmt.Println(groupName)
	selfPtr.messageQueue.Enlist(groupName, msg)

	selfPtr.waiterCountLock.Lock()
	if *selfPtr.waiterCount > 0 {
		selfPtr.messageArrivalChannel <- 1
	}
	selfPtr.waiterCountLock.Unlock()


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
		selfPtr.waiterCountLock.Lock()
		*selfPtr.waiterCount = *selfPtr.waiterCount + 1
		selfPtr.waiterCountLock.Unlock()

		<-selfPtr.messageArrivalChannel

		selfPtr.waiterCountLock.Lock()
		*selfPtr.waiterCount -= 1
		selfPtr.waiterCountLock.Unlock()

		message = selfPtr.availableMessage()
//		fmt.Println("got the message : ")
//		fmt.Println(message)
		return message
	}
	return message
}

func (selfPtr *PriorityBasedChannel) setPriorityForMessageType(typ reflect.Type, priority int) {
	typeName := getTypeNameFromType(typ)
	selfPtr.messageQueue.SetGroup(typeName, priority)
	(*selfPtr.priorityMap)[typeName] = priority
}
