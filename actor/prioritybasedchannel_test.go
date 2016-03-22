package actor
import (
	"testing"
//	"fmt"
	"sync"
	"reflect"
)

var count int = 0
var mutex sync.Mutex

func cuser2(c PriorityBasedChannel, quitChannel chan int) {

	c.Get()
	quitChannel <- 1
	mutex.Lock()
	count = count + 1

	mutex.Unlock()
}

type Msg struct {
	name string
	actoRef ActorRef
}

func TestSome(t *testing.T) {
	c := NewPriorityBasedChannel("TestChannel")

	c.SetPriority(1, reflect.TypeOf(Msg{}))
	quitChannel := make(chan int)

	count := 10000
	for i := 0; i < count; i++ {
		go cuser2(c, quitChannel)
	}

	for i := 0; i < count; i++ {
		c.Send(Msg{name:"yo"})
	}

	for i := 0; i < count; i++ {
		<- quitChannel
	}
}
