package actor

/**
A new actor must implement this to be created
 */
type ActorInterface interface {
	OnReceive(self ActorRef, msg ActorMessage) error
	OnStart(self ActorRef) error
	OnStop(self ActorRef) error
	OnRestart(self ActorRef) error
	OnDeadletter(self ActorRef) error
	OnError(self ActorRef)
}


type DefaultActorInterface struct {

}

func (selfPtr *DefaultActorInterface) OnReceive(self ActorRef, msg ActorMessage) error {
	return nil
}

func (selfPtr *DefaultActorInterface) OnStart(self ActorRef) error {
	return nil
}

func (selfPtr *DefaultActorInterface) OnStop(self ActorRef) error {
	return nil
}

func (selfPtr *DefaultActorInterface) OnRestart(self ActorRef) error {
	return nil
}

func (selfPtr *DefaultActorInterface) OnDeadletter(self ActorRef) error {
	return nil
}

func (selfPtr *DefaultActorInterface) OnError(self ActorRef) {

}