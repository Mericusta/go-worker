package fsm

type FSMState int

var WaittingState chan bool
var ExecuteState chan bool
var ExitState chan bool

func init() {
	WaittingState = make(chan bool)
	ExecuteState = make(chan bool)
	ExitState = make(chan bool)
}

func FSM() {
	// select
}
