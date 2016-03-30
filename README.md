# chucktor
Chucktor : Actor implementation with golang channels ! 

Chucktor is a light weight actor model implementation for go language.

## Why actors when channels ##

Channels are great. However actors provide a hierarchical structure which sometimes makes a lot of sense, depending on the project. Underneath this implementation there are a lot of channels anyway :).

## How to use ##

The expected usage is to create one actor system per application. An actor system can be created with the statement:

```go
executorCount := 4
actorSystem := actor.NewActorSystem(executorCount)
```

Executor count represents the alive goroutine workers that are responsible to process messages sent to actors.
An actor is created either via an actor system object or an actor reference.

An actor is declared as a struct that implements actor.ActorInterface. ActorInterface is declared as follows:

```go
type ActorInterface interface {
	OnReceive(self ActorRef, msg ActorMessage) error
	OnStart(self ActorRef) error
	OnStop(self ActorRef) error
	OnRestart(self ActorRef) error
	OnDeadletter(self ActorRef) error
}
```

An implementor of this interface should use actor.DefaultActorInterface via composition to be able to use default implementations most of which do nothing. But yet it prevents one having to implement the unused methods.

An example that creates an actor system and actor on this system, can be seen below:

```go
import "actor"

//composition for default callback implementations
type MyActor struct {
  actor.DefaultActorInterface
}

type MessageStruct struct {
  text string
}

//On receive callback implementation
func (selfPtr *MyActor) OnReceive(self actor.ActorRef, msg actor.ActorMessage) error {
  message := msg.Msg
  switch message.(type) {
    case MessageStruct:
      fmt.Println("I have received a message of the expected type")
    default:
      fmt.Println("Message is not of the expected type")
      
  return nil
  }
}

func main() {
  executorCount := 4
  actorSystem := actor.NewActorSystem(executorCount)
  var actor1ref actor.ActorRef := actorSystem.CreateActor(&MyActor{}, "actor1")
  
  //Send a message to actor1 from a nil sender
  actor1ref.Tell(MessageStruct{}, nil)
  
  actor.Run()
}
```

