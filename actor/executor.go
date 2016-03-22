package actor
import "fmt"


type ActorSystemExecutor struct {
	actorSystem *ActorSystem
}

//Execute a message from the grabbed actor from the queue channel
func (selfPtr *ActorSystemExecutor) execute(actor *CoreActor) {
	fmt.Println("EXECUTING")
	actor.run()
	fmt.Println("DONE EXECUTING")
}

//Starts eternal execution
func (selfPtr *ActorSystemExecutor) startExecution() {
	for {
		defaultActor := <- selfPtr.actorSystem.actorChannel
		selfPtr.execute(defaultActor)
	}
}

func newActorSystemExecutor(actorSystem *ActorSystem) *ActorSystemExecutor {
	return &ActorSystemExecutor{actorSystem:actorSystem}
}
