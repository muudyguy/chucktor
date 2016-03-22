package actor
import (
	"os"
	"os/signal"
	"fmt"
	"reflect"
)

func convertDefaultActorToActorRef(defaultActor *CoreActor) ActorRef {
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

func getTypeNameFromType(typ reflect.Type) string {
	nameOfType := typ.PkgPath() + typ.Name()
	return nameOfType
}


func getGroupNameFromMsg(msg interface{}) string {
	typ := reflect.TypeOf(msg)
	return getTypeNameFromType(typ)
}