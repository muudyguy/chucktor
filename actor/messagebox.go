package actor

type MessageBoxInterface interface {
	Pop() *CoreActor
	Enlist(coreActor *CoreActor)
}


