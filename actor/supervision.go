package actor

type SupervisionConstants struct {
	PROPAGATE int
	RESTART int
}


var SUPERVISION SupervisionConstants = SupervisionConstants{PROPAGATE:0, RESTART:1}
