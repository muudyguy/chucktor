package actor
import (
	"reflect"
	"sync"
	"queue"
	"fmt"

	"sync/atomic"
	"strconv"
)

type ChannelInterface interface {
	Send(msg interface{})
	Get() interface{}
}

type PriorityBasedChannel struct {
	channelName           string

	messageArrivalChannel chan int

	priorityMap           *map[string]int
	messageQueue          *queue.RoundRobinQueue

	waiterCount           *uint32

	priorityMapLock       *sync.Mutex
	enlistGetOneLock      *sync.Mutex
}

func NewPriorityBasedChannel(name string) PriorityBasedChannel {
	pbc := PriorityBasedChannel{}
	pbc.channelName = name

	pbc.priorityMap = new(map[string]int)
	*pbc.priorityMap = make(map[string]int)

	pbc.priorityMapLock = new(sync.Mutex)
	pbc.enlistGetOneLock = new(sync.Mutex)

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

	fmt.Println("Enlisted item now count of channel " + selfPtr.channelName + " is " + strconv.Itoa(selfPtr.messageQueue.GetTotalItemCount()))

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


	GetStart:
	message := selfPtr.availableMessage()
	//If message was nil, there are no available messages at the moment
	//So start waiting
	if message == nil {
		//increase waiter count
		atomic.AddUint32(selfPtr.waiterCount, 1)
		//start waiting
		<-selfPtr.messageArrivalChannel
		//reduce waiter count
		atomic.AddUint32(selfPtr.waiterCount, ^uint32(0))

		//Now that we received a message from enlist method, Get the message
		message = selfPtr.availableMessage()

		/**
			If the message is nil again, there are some possibilities
			1 - Enlist Sends => Get receives message but does not yet reduce waiter count,
			 	another Enlist comes, sees that waiter count is bigger than 0, sends again. Now there are two messages
			 	Get receives the first message, and comes to GetStart label, receives again, skips this if block
			 	comes to GetStart again. Now there is 1 message in messageArrivalChannel, but 0 messages in the queue.
			 	Get receives nil from first availableMessage method, then gets in this if block. Gets the arrival message
			 	from channel, then again receives nil from available message. If now we return this from Get method, the
			 	receiver will break. So we need to go back to GetStart label, for all these false Gets.
			 	todo Find a better way ??
		 */
		if message == nil {
			goto GetStart
		}
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
