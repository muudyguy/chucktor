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
	"strconv"
//	"runtime"
)

type Msg struct {
	name string
	actoRef actor.ActorRef
}

type Msg2 struct {
	name string
}

type MyActor struct {
	actor.DefaultActorInterface
}

func (selfPtr *MyActor) OnReceive(self actor.ActorRef, msg actor.ActorMessage) error {
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

	return nil
}

func (selfPtr *MyActor) OnStart(self actor.ActorRef) error {
	fmt.Println("onStart callback Ran")
	return nil
}

func (selfPtr *MyActor) OnStop(self actor.ActorRef) error {
	fmt.Println("onStop callback ran")
	return nil
}

func (selfPtr *MyActor) OnRestart(self actor.ActorRef) error {
	fmt.Println("onRestart callback ran")
	return nil
}

type MyActor2 struct {
	actor.DefaultActorInterface
}

func (myActor MyActor2) OnReceive(self actor.ActorRef, msg actor.ActorMessage) error {
	switch msg.Msg.(type) {
	case Msg:
		teller := msg.Teller
		teller.Tell(Msg2{"answer"}, self)
	default:
		fmt.Println("dont know the message")
	}
	return nil
}

func (selfPtr *MyActor2) OnStart(self actor.ActorRef) error {
	return nil
}

func (selfPtr *MyActor2) OnStop(self actor.ActorRef) error {
	return nil
}

func (selfPtr *MyActor2) OnRestart(self actor.ActorRef) error {
	return nil
}

func testActorMessaging() {
	as := actor.NewActorSystem(4)
	myActor := MyActor{}
	myActor2 := MyActor2{}
	actor1, err := as.CreateActor(&myActor, "trip")
	actor1.SetPriority(reflect.TypeOf(Msg{}), 1)
	actor1.SetPriority(reflect.TypeOf(Msg2{}), 1)
	if err != nil {
		panic(err)
	}


	actor2, err2 := as.CreateActor(&myActor2, "trip/trip2")
	actor2.SetPriority(reflect.TypeOf(Msg{}), 1)
	actor2.SetPriority(reflect.TypeOf(Msg2{}), 1)
	if err2 != nil {
		panic(err)
	}

	actor1.Tell(Msg{name:"name", actoRef:actor2}, actor.ActorRef{})

	for i := 0; i < 2; i++ {
		actor1.Tell(Msg{name:"name", actoRef:actor2}, actor.ActorRef{})
//		actor1.Tell(Msg{name:"name", actoRef:actor2}, actor.ActorRef{})
	}
}

func testActorCreation() {
	as := actor.ActorSystem{}
	as.InitSystem()

	for i := 0; i < 1000000; i++ {
		myActor := MyActor{}
		actor1, err := as.CreateActor(&myActor, strconv.Itoa(i))
		if err != nil {
			panic(err)
		}
		fmt.Println(actor1)

	}
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
//	fmt.Println(runtime.GOMAXPROCS(0))
	testActorMessaging()

	actor.Run()

}
