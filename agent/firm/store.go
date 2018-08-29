package firm

import (
	"log"

	"github.com/ninjadotorg/economy-simulation/agent"
	"github.com/ninjadotorg/economy-simulation/state"
)

type Store struct{}

func (st *Store) Init(a *agent.Agent) {
	log.Println("init store firm agent")
}

func (st *Store) Run(a *agent.Agent, s state.State) {
	log.Println("i'm a store")
}
