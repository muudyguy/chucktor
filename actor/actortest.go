package actor

import (
	"time"
	"testing"
	"reflect"
	"fmt"
)

type TestActor struct {
	DefaultActorInterface
	messageReceivalChannel chan interface{}
}

func NewTestActor(channel chan interface{}) *TestActor {
	ta := TestActor{}
	ta.messageReceivalChannel = channel
	return &ta
}

func (selfPtr *TestActor) OnReceive(self ActorRef, msg ActorMessage) error {
	fmt.Println("GOT THE MESSAGE BABE")
	selfPtr.messageReceivalChannel <- msg.Msg
	return nil
}


type ActorTestContext struct {
	testActor ActorRef
	channel chan interface{}
	t *testing.T
	actorSystem *ActorSystem
}

func (selfPtr *ActorTestContext) Tell(receiver ActorRef, msg interface{}) {
	receiver.Tell(msg, selfPtr.testActor)
}

type UnexpectedMessageType struct {

}

func (self UnexpectedMessageType) Error() string {
	return "Unexpected message type"
}

type Timeout struct {

}

func (self Timeout) Error() string {
	return "Timeout at expectMsg"
}

func (selfPtr *ActorTestContext) expectMsg(seconds int, msg interface{}) error {
	//todo Receive from channel with timeout ?
	timeoutChannel := make(chan int)
	go func() {
		time.Sleep(time.Duration(seconds) * time.Second)
		timeoutChannel <- 1
	}()

	select {
	case <- timeoutChannel:
		return Timeout{}
	case receivedMsg := <- selfPtr.channel:
		//According to the spec interfaces can be compared in a viable way
		if receivedMsg != msg {
			return UnexpectedMessageType{}
		}
	}

	return nil
}

func (selfPtr *ActorTestContext) ExpectMsg(seconds int, msg interface{}) {
	err := selfPtr.expectMsg(seconds, msg)
	if err != nil {
		selfPtr.t.Error(err)
	}
}

func NewActorTestContext(t *testing.T, actorSystem *ActorSystem) *ActorTestContext {
	atc := ActorTestContext{}
	atc.t = t
	atc.actorSystem = actorSystem
	channel := make(chan interface{})
	atc.channel = channel
	testActorRef, err := actorSystem.CreateActor(NewTestActor(channel), "testactor")
	if err != nil {
		panic(err)
	}
	testActorRef.SetPriority(reflect.TypeOf(5), 1)
	atc.testActor = testActorRef
	return &atc
}