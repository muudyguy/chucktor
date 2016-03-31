# chucktor
Chucktor : Actor implementation with golang channels ! 

Chucktor is a light weight actor model implementation for go language.

# Why actors when channels #

Channels are great. However actors provide a hierarchical structure which sometimes makes a lot of sense, depending on the project. Underneath this implementation there are a lot of channels anyway :).

# How to use #

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

## What can you do with an actor reference? ##

1- Tell a message to it
2- Stop it now
3- Stop it when it is done with current messages
3- Restart it

### Telling an actor a message ###
```go
actorRef.Tell(msg, tellerActorRef)
```

When a message is to be told to another actor, the reference of the teller needs to be passed into the Tell method.
That is the only way the receiver can know the sender of the message.

### Stopping an Actor ###

```go
actorRef.Stop()
```

Currently, when an actor is stopped, the stopped actor cannot know the identity of the stopper. This will be fixed in later releases.

### Stopping an actor when done ###

To be implemented... Very soon :)

### Restart an actor ###

Restarting an actor simply just, invokes the onRestart method. Currently old messages are not deleted or processed as dead letter.

```go
actorRef.Restart()
```

In the future, restarter will have to be passed into Restart method.

## Message Priority in Actors ##

Messages in chucktor are processed with a scheduler very similar to round-robin. Every message type has a quantum that it can use. Once quantum of a message type is completed, next type is started to be processed. Message priority can be set as follows:

```go
var typeDummy int = 5
actorRef.SetPriority(reflect.TypeOf(typeDummy), 1)
```

typeDummy variable is used to create a Type variable. This is obviously not a very clean way to it, but is the only way to set priorities currently.

1 is the priority number. System has an internal quantum setting, that currently cannot be edited from outside. 
Assuming we have ran:

```go
var typeDummyInt int = 5
var typeDummyString string = "a"
actorRef.SetPriority(reflect.TypeOf(typeDummyInt), 1)
actorRef.SetPriority(reflect.TypeOf(typeDummyString), 2)
```

and quantum is 5, and we have 10 integer messages to 20 string messages, actor will process in this order:

10 string -> 5 int -> 10 string -> 5 int

Basically, (quantum * priority) is the maximum times a message type will be processed until next type is promoted.
If there are not enough message for this maximum number, next message type will be promoted prematurely.

If there is no priority set for an actor reference, every message will have the same priority. 


