package economy

//type State map[string]interface{}

type State struct {
	price  float64
	block  int
	agents []Agent
}

func currentState(econ Economy) State {
	s := State{}
	s.price = 1
	s.block = 1000
	return s
}
