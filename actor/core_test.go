package actor
import (
	"testing"
	"reflect"
	"fmt"
	"time"
)


/**
 **
 **
 **
Tests an actor sending message to another actor with multiple executors ( 4 )
 **
 **
 **
 */
type MyTestActor struct {
	DefaultActorInterface
}

func (selfPtr *MyTestActor) OnReceive(self ActorRef, msg ActorMessage) error {
	return nil
}

type MyTestActor2 struct {
	DefaultActorInterface
	quitChannel chan int
	t *testing.T
	testActor1ref ActorRef
}

func (selfPtr *MyTestActor2) OnReceive(self ActorRef, msg ActorMessage) error {
	message := msg.Msg.(int)
	if message != 5 {
		selfPtr.t.Error("Message is incorrect")
	}

	if selfPtr.testActor1ref != msg.Teller {
		selfPtr.t.Error("Sender actor ref is incorrect")
	}
	selfPtr.quitChannel <- 1
	return nil
}

func tellMessages(testActor1ref ActorRef, testActor2ref ActorRef, count int) {
	for i := 0; i < count; i++ {
		testActor2ref.Tell(5, testActor1ref)
	}
}

func TestMessaging(t *testing.T) {
	quitChannel := make(chan int)
	actorSystem := NewActorSystem(4)
	testActor1 := MyTestActor{}


	testActor1ref, err := actorSystem.CreateActor(&testActor1, "testactor1")
	if err != nil {
		panic(err)
	}

	testActor2 := MyTestActor2{quitChannel:quitChannel, t:t, testActor1ref:testActor1ref}
	testActor2ref, err := actorSystem.CreateActor(&testActor2, "testactor1")
	if err != nil {
		panic(err)
	}

	testActor1ref.SetPriority(reflect.TypeOf(5), 1)
	testActor2ref.SetPriority(reflect.TypeOf(5), 1)

	count := 1000
	go tellMessages(testActor1ref, testActor2ref, count)

	for i := 0; i < count; i++ {
		<- quitChannel
	}

}


/**
*
*
*
THIS BLOCK TESTS THAT STOPPING ACTOR WORKS !
*
*
*
 */
type StopTestActor struct {
	DefaultActorInterface
	onStopChannel chan int
	messageCount int
	maximumMessageCount int
	t *testing.T
}

func (selfPtr *StopTestActor) OnReceive(self ActorRef, msg ActorMessage) error {
	fmt.Println("On receive of stop test actor")
	selfPtr.messageCount += 1
	if selfPtr.messageCount == selfPtr.maximumMessageCount {
		selfPtr.t.Error("All messages were processed via on receive")
	}
	return nil
}

func (selfPtr *StopTestActor) OnStop(self ActorRef) error {
	fmt.Println("On stop of stop test actor")
	selfPtr.onStopChannel <- 1
	return nil
}



func sendMessagesToStopActor(stopActorRef ActorRef, count int, t *testing.T, done chan int) {
	for i := 0; i < count; i++ {
		err := stopActorRef.Tell(5, ActorRef{})
		if err != nil {
			fmt.Println("Waiting at done")
			done <- 1
			return
		}
	}
	done <- 1
}

func TestStopping(t *testing.T) {
	onStopChannel := make(chan int)
	actorSystem := NewActorSystem(4)

	count := 1000
	stopTestActor1 := StopTestActor{onStopChannel:onStopChannel, maximumMessageCount:count, t:t}


	testActor1ref, err := actorSystem.CreateActor(&stopTestActor1, "testactor1")
	if err != nil {
		panic(err)
	}

	testActor1ref.SetPriority(reflect.TypeOf(5), 1)

	done := make(chan int)


	go sendMessagesToStopActor(testActor1ref, count, t, done)
	//Wait 1 second for purposes
	time.Sleep(100000)
	testActor1ref.Stop()
	<- done
	select {
	case <- onStopChannel:
		//no problem
	case <- time.After(time.Second * 3):
		t.Error("On stop was not entered for 3 seconds")
	default:
		t.Error("On stop was not entered")
	}
}