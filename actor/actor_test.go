package actor
import (
	"testing"
	"time"
	"fmt"
)

type WatcherActor struct {
	DefaultActorInterface
	doneChannel chan int
}

func (selfPtr *WatcherActor) OnReceive(self ActorRef, msg ActorMessage) error {
	switch msg.Msg.(type) {
	case Death:
		fmt.Println("Received death")
		selfPtr.doneChannel <- 1
	default:

	}
	return nil
}


type WatchedActor struct {
	DefaultActorInterface
}

func (selfPtr *WatchedActor) OnReceive(self ActorRef, msg ActorMessage) error {
	return nil
}


func TestWatch(t *testing.T) {
	as := NewActorSystem(4)
	doneChannel := make(chan int)
	watcher1ref, err := as.CreateActor(&WatcherActor{doneChannel:doneChannel}, "watcher1")
	if err != nil {
		panic(err)
	}
	watcher2ref, err := as.CreateActor(&WatcherActor{doneChannel:doneChannel}, "watcher2")
	if err != nil {
		panic(err)
	}
	watchedref, err := as.CreateActor(&WatchedActor{}, "watched")

	watchedref.Watch(watcher1ref)
	watchedref.Watch(watcher2ref)

	time.Sleep(time.Duration(1) * time.Second)
	watchedref.Stop()

	<- doneChannel
	<- doneChannel
}


/**
FAIL TEST
 */

type FailingActor struct {
	DefaultActorInterface
}

func (selfPtr *FailingActor) OnReceive(self ActorRef, msg ActorMessage) error {
	fmt.Println("IN FAILINGACTOR ON RECEIVE")
	return fmt.Errorf("Error")
}

type ParentOfFailingActor struct {
	DefaultActorInterface
	failReceivedChannel chan int
}

func (selfPtr *ParentOfFailingActor) OnReceive(self ActorRef, msg ActorMessage) error {
	return nil
}

func (selfPtr *ParentOfFailingActor) OnError(self ActorRef) {
	selfPtr.failReceivedChannel <- 1
}

func TestFail(t *testing.T) {
	as := NewActorSystem(4)
	failReceivedChannel := make(chan int)
	_, err := as.CreateActor(&ParentOfFailingActor{failReceivedChannel:failReceivedChannel}, "parentoffailingactor")
	if err != nil {
		panic(err)
	}
	failingActor, err := as.CreateActor(&FailingActor{}, "parentoffailingactor/failingactor")
	if err != nil {
		panic(err)
	}
	failingActor.Tell(5, ActorRef{})
	<- failReceivedChannel
}