package actor

type ActorNotFound struct {

}

func (actorNotFound ActorNotFound) Error() string {
	return "Actor for Path Not Found"
}

type ActorIsStopped struct {

}

func (actorNotFound ActorIsStopped) Error() string {
	return "Actor is stopped"
}
