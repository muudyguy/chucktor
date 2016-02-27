package actor
import (
	"testing"
//	"fmt"
	"sync"
	"reflect"
)

var count int = 0
var mutex sync.Mutex

func cuser2(c PriorityBasedChannel) {

	//		fmt.Println("gonna get now")
	c.Get()
//	fmt.Println(c.messageQueue.GetTotalItemCount())

	//		fmt.Println(m)
	mutex.Lock()
	count = count + 1

	mutex.Unlock()
}

type Msg struct {
	name string
	actoRef ActorRef
}

func TestSome(t *testing.T) {
	c := NewPriorityBasedChannel()

	c.SetPriority(1, reflect.TypeOf(Msg{}))

	for i := 0; i < 10000; i++ {
		go cuser2(c)
	}

	for i := 0; i < 10000; i++ {
		c.Send(Msg{name:"yo"})
	}

	Run()
}
