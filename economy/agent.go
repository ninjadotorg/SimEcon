package economy

type Agent struct {
	agentType string
	balanceOf map[string]float64 // map[asset]size
}

func (agent *Agent) balance(asset string) float64 {
	return agent.balanceOf[asset]
}
