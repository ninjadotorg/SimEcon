package sim

type State map[string]interface{}

func state() State {
	s := State{}
	s["price"] = 1
	s["block"] = 1000
	return s
}
