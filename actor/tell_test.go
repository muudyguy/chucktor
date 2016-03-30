package actor
import (
	"testing"
	"reflect"
)


/**
SIMPLE TELL TEST
 */
type SenderActor struct {
	DefaultActorInterface
}

func (selfPtr *SenderActor) OnReceive(self ActorRef, msg ActorMessage) error {
	return nil
}

type ReceiverActor struct {
	DefaultActorInterface
	receivedChannel chan int
	senderRef       ActorRef
	t	*testing.T
}

func (selfPtr *ReceiverActor) OnReceive(self ActorRef, msg ActorMessage) error {
	if msg.Teller != selfPtr.senderRef {
		selfPtr.t.Error("Sender is not supposed to be who it is")
	}

	if msg.Msg.(int) != 5 {
		selfPtr.t.Error("Message is incorrect")
	}

	selfPtr.receivedChannel <- 1
	return nil
}


func TestSimpleTell(t *testing.T) {
	as := NewActorSystem(4)
	receivedChannel := make(chan int)


	sender, err := as.CreateActor(&SenderActor{}, "sender")
	if err != nil {
		panic(err)
	}
	sender.SetPriority(reflect.TypeOf(5), 1)

	receiver, err := as.CreateActor(&ReceiverActor{receivedChannel:receivedChannel,
			senderRef:sender,
			t:t}, "receiver")

	if err != nil {

	}
	receiver.SetPriority(reflect.TypeOf(5), 1)

	err = receiver.Tell(5, sender)
	if err != nil {
		panic(err)
	}

	<- receivedChannel

}


/**
Test tell parent
 */
func TestTellParent(t *testing.T) {
	as := NewActorSystem(4)
	receivedChannel := make(chan int)

	receiverActorPtr := &ReceiverActor{receivedChannel:receivedChannel,
		t:t}
	receiverRef, err := as.CreateActor(receiverActorPtr, "receiver")

	if err != nil {
		panic(err)
	}
	senderRef, err := as.CreateActor(&SenderActor{}, "receiver/sender")

	receiverActorPtr.senderRef = senderRef

	senderRef.SetPriority(reflect.TypeOf(5), 1)
	receiverRef.SetPriority(reflect.TypeOf(5), 1)

	senderRef.TellParent(5, senderRef)

	<- receivedChannel
}


/**
Test tell children
 */
func TestTellChildren(t *testing.T) {
	as := NewActorSystem(4)
	receivedChannel := make(chan int)

	senderRef, err := as.CreateActor(&SenderActor{}, "sender")
	if err != nil {
		panic(err)
	}

	receiverActorPtr := &ReceiverActor{receivedChannel:receivedChannel, senderRef: senderRef, t:t}
	_, err = as.CreateActor(receiverActorPtr, "sender/receiver1")
	if err != nil {
		panic(err)
	}
	_, err = as.CreateActor(receiverActorPtr, "sender/receiver2")
	if err != nil {
		panic(err)
	}

//	senderRef.SetPriority(reflect.TypeOf(5), 1)
//	receiverRef1.SetPriority(reflect.TypeOf(5), 1)
//	receiverRef2.SetPriority(reflect.TypeOf(5), 1)

	senderRef.TellChildren(5, senderRef)

	<- receivedChannel
	<- receivedChannel
}