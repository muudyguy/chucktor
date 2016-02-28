package main

import (
	"muddle/actor"
	"fmt"
//	"muddle/tree"
//	"strings"
//	"queue"
//	"reflect"
	"reflect"
	"sync"
//	"os"
)

type Msg struct {
	name string
	actoRef actor.ActorRef
}

type Msg2 struct {
	name string
}

type MyActor struct {

}

func (selfPtr *MyActor) OnReceive(self actor.ActorRef, msg actor.ActorMessage) {
	switch msg.Msg.(type) {
	case Msg:
		message := msg.Msg.(Msg)
		actorToForward := message.actoRef
		actorToForward.Tell(message, self)
	case Msg2:
//		fmt.Println("I am actor1 and received reply from actor2, it has in it : " + msg.Msg.(Msg2).name)
	default:
//		fmt.Println("Dont know")
	}
}

func (selfPtr *MyActor) OnStart(self actor.ActorRef) {

}

func (selfPtr *MyActor) OnStop(self actor.ActorRef) {

}

type MyActor2 struct {

}

func (myActor MyActor2) OnReceive(self actor.ActorRef, msg actor.ActorMessage) {
	switch msg.Msg.(type) {
	case Msg:
//		fmt.Println("received message from actor1")
		teller := msg.Teller
		teller.Tell(Msg2{"answer"}, self)
	default:
		fmt.Println("dont know the message")
	}
}

func (selfPtr *MyActor2) OnStart(self actor.ActorRef) {

}

func (selfPtr *MyActor2) OnStop(self actor.ActorRef) {

}
//
//func (msg Msg) Bigger(in tree.Comparable) bool {
//	toBeCompared := in.(*Msg)
//	return msg.name > toBeCompared.name
//}
//
//func (msg Msg) Equals(in tree.Comparable) bool {
//	toBeCompared := in.(*Msg)
//	return msg.name == toBeCompared.name
//}

func testActorMessaging() {
	as := actor.ActorSystem{}
	as.InitSystem()
	myActor := MyActor{}
	myActor2 := MyActor2{}
	actor1, err := as.CreateActor(&myActor, "trip")
	if err != nil {
		panic(err)
	}


	actor2, err2 := as.CreateActor(&myActor2, "trip2")
	if err2 != nil {
		panic(err)
	}

	actor1.Tell(Msg{name:"name", actoRef:actor2}, actor.ActorRef{})

	if err != nil {
		panic(err)
	}
	actor1.Tell(Msg{name:"name", actoRef:actor2}, actor.ActorRef{})

	actor.Run()
}

var count int = 0
var mutex sync.Mutex

func cuser2(c actor.PriorityBasedChannel) {

//		fmt.Println("gonna get now")
		c.Get()


		mutex.Lock()
		count = count + 1
//		if count > 101 {
//			os.Exit(0)
//		}


		fmt.Println(count)




		mutex.Unlock()
}

func testChannel() {
	c := actor.NewPriorityBasedChannel("c")

	c.SetPriority(1, reflect.TypeOf(Msg{}))

	//	for i := 0; i < 10000; i++ {
	//		go cuser2(c)
	//	}
	//
	//	for i := 0; i < 10000; i++ {
	//		c.Send(Msg{name:"yo"})
	//	}

	for i := 0; i < 10; i++ {
		c.Send(Msg{name:"yo"})
		go cuser2(c)
	}
}

func testChannelGet() {
	c := actor.NewPriorityBasedChannel("c")

	c.SetPriority(1, reflect.TypeOf(Msg{}))

	c.Send(Msg{name:"yo"})
	c.Send(Msg{name:"y2"})

	fmt.Println(c.Get())

	fmt.Println(c.Get())
}

func main() {

	testActorMessaging()

	actor.Run()

}
