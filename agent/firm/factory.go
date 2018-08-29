package firm

import (
	"log"

	"github.com/ninjadotorg/economy-simulation/agent"
	"github.com/ninjadotorg/economy-simulation/state"
)

type Factory struct {
	salaries map[string]float64 // map(agent id, salary)
}

func (f *Factory) Init(a *agent.Agent) {
	log.Println("init factory")
}

func (f *Factory) Run(a *agent.Agent, s state.State) {
	log.Println("i'm a factory")
}
