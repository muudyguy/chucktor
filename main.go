package main

import (
	"muddle/actor"
	"fmt"
	"muddle/tree"
)

type Msg struct {
	name string
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

func (msg Msg) Bigger(in tree.Comparable) bool {
	toBeCompared := in.(*Msg)
	return msg.name > toBeCompared.name
}

func (msg Msg) Equals(in tree.Comparable) bool {
	toBeCompared := in.(*Msg)
	return msg.name == toBeCompared.name
}

func main() {
//	as := actor.ActorSystem{}
//	as.InitSystem()
//	myActor := MyActor{}
//	actor := as.CreateActor(myActor, "actorname")
//	actor.Tell(Msg{}, nil)


	var btree tree.BinaryTree = tree.NewBinaryTree()
	btree.Add(&Msg{"selam"})

	a := btree.Search(&Msg{"selam"})
	fmt.Printf("%p", a)
}
