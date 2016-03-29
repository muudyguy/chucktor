# chucktor
Chucktor : Actor implementation with golang channels ! 

Chucktor is a light weight actor model implementation for go language.

## Why actors when channels ##

Channels are great. However actors provide a hierarchical structure which sometimes makes a lot of sense, depending on the project. Underneath this implementation there are a lot of channels anyway :).

## How to use ##

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

