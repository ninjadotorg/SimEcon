package state

//type State map[string]interface{}

type State struct {
	Price float64
	Block int
}

func CurrentState() State {
	s := State{}
	s.Price = 1
	s.Block = 1000
	return s
}
