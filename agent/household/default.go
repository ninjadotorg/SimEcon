package household

import (
	"log"

	"github.com/ninjadotorg/economy-simulation/agent"
	"github.com/ninjadotorg/economy-simulation/state"
)

type Default struct {
	employed bool
}

func (d *Default) Init(a *agent.Agent) {
	log.Println("init default household agent")
}

func (d *Default) Run(a *agent.Agent, s state.State) {
	log.Println("i'm a default household agent")
	if !d.employed {
		// look for a job
	}
}
