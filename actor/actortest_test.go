package actor
import (
	"testing"
	"fmt"
)

type MyActor struct{
	DefaultActorInterface
}

type MyMsg struct {
	Text string
}

func (selfPtr *MyActor)OnReceive(self ActorRef, msg ActorMessage) error {
	msg.Teller.Tell(5, self)
	return nil
}

//Tests whether a simple expectMsg operation works
func TestExpectMsg(t *testing.T) {
	as := NewActorSystem(4)
	actorRef, err := as.CreateActor(&MyActor{}, "testactor")
	if err != nil {
		panic(err)
	}
	testContext := NewActorTestContext(t, as)
	testContext.Tell(actorRef, 1)

	err = testContext.expectMsg(5, 5)
	if err != nil {
		t.Error("ExpectMsg was not supposed to throw an error !")
	}
}


func TestExpectMsgFailWhenNoMessageReceivedWithExpectedValue(t *testing.T) {
	defer func() {
		fmt.Println(recover())
		if r := recover(); r != nil {
			fmt.Println(r)
		}
	}()
	as := NewActorSystem(4)
	actorRef, err := as.CreateActor(&MyActor{}, "testactor")
	if err != nil {
		panic(err)
	}
	testContext := NewActorTestContext(t, as)
	testContext.Tell(actorRef, 1)

	//Expecting 2, but 5 will be received
	err = testContext.expectMsg(5, 2)
	if err == nil {
		t.Error("There was expected an unexpected message type error, but threw no error")
	}
}
