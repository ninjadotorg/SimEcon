package sim

import (
	"log"
)

type AgentGroup struct {
	Type     string `json:"type"`
	Qty      int    `json:"qty"`
	StepSize int    `json:"stepSize"`
}

type Agent struct {
	name     string
	stepSize int
	state    chan State
	quit     chan struct{}
}

func newAgent(agentType string, stepSize int) (a Agent) {
	a.stepSize = stepSize
	a.name = agentType
	a.state = make(chan State, 10)
	return
}

func (a *Agent) start() {
	for {
		state := <-a.state
		log.Println("agent:", a.name, "state:", state)
	}
}
