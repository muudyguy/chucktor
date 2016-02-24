package actor
import (
	"os"
	"os/signal"
	"fmt"
)

func convertDefaultActorToActorRef(defaultActor *DefaultActor) ActorRef {
	actorRef := ActorRef{defaultActor.index, defaultActor}
	return actorRef
}

func Run() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)
	<-c
	fmt.Println("Exit signal received from os !")
	os.Exit(0)
}