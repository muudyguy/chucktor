package actor

type ActorNotFound struct {

}

func (actorNotFound ActorNotFound) Error() string {
	return "Actor for Path Not Found"
}
