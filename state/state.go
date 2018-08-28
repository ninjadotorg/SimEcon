package state

type State map[string]interface{}

func CurrentState() State {
	s := State{}
	s["price"] = 1
	s["block"] = 1000
	return s
}
