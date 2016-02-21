package main

import (
	"muddle/actor"
	"fmt"
)

type Msg struct {

}

type MyActor struct {

}

func (myActor MyActor) OnRecieve(actor *actor.ActorRef, msg actor.ActorMessage) {
	switch msg.Msg.(type) {
	case Msg:
		fmt.Println("Hell yea")
	default:
		fmt.Println("Dont know")
	}
}

func main() {
	as := actor.ActorSystem{}
	as.InitSystem()
	myActor := MyActor{}
	actor := as.CreateActor(myActor, "actorname")
	actor.Tell(Msg{}, nil)
}
