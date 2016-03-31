package actor



type ActorSystemExecutor struct {
	actorSystem *ActorSystem
}

//Execute a message from the grabbed actor from the queue channel
func (selfPtr *ActorSystemExecutor) execute() {
	coreActor := selfPtr.actorSystem.messageBox.Pop()
	coreActor.run()
}

//Starts eternal execution
func (selfPtr *ActorSystemExecutor) startExecution() {
	for {
		<- selfPtr.actorSystem.actorChannel
		selfPtr.execute()
	}
}

func newActorSystemExecutor(actorSystem *ActorSystem) *ActorSystemExecutor {
	return &ActorSystemExecutor{actorSystem:actorSystem}
}
