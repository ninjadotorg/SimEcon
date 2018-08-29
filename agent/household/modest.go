package household

import (
	"log"

	"github.com/ninjadotorg/economy-simulation/agent"
	"github.com/ninjadotorg/economy-simulation/state"
)

type Modest struct {
	// state
}

func (m *Modest) Init(a *agent.Agent) {
	log.Println("init modest household agent")
}

func (m *Modest) Run(a *agent.Agent, s state.State) {
	log.Println("i'm a modest household")
}
