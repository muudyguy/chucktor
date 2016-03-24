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

func (selfPtr *MyActor) OnReceive(self actor.ActorRef, msg actor.ActorMessage) error {
  message := msg.Msg
  switch message.(type) {
    
  }
}

func main() {
  executorCount := 4
  actorSystem := actor.NewActorSystem(executorCount)
  actorSystem.CreateActor(
}
```

