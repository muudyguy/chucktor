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