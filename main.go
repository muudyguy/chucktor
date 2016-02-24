package main

import (
	"muddle/actor"
	"fmt"
	"muddle/tree"
	"strings"
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
	s:= "/a/b"
	fmt.Println(len(strings.Split(s, "/")))

	as := actor.ActorSystem{}
	as.InitSystem()
	myActor := MyActor{}
	actor, err := as.CreateActor(myActor, "trip")
	if err != nil {
		fmt.Println(err)
		return
	}
	actor, err = as.CreateActor(myActor, "/trip/actorname")
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(actor)
	actor.Tell(Msg{}, nil)

	ar, err := as.GetActorRef("trip/actorname")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(ar)
	var btree tree.BinaryTree = tree.NewBinaryTree()
	btree.Add(&Msg{"selam"})

//	a := btree.Search(&Msg{"selam"})
//	fmt.Printf("%p", a)


}
