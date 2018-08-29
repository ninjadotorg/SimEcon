package household

import (
	"log"

	"github.com/ninjadotorg/economy-simulation/agent"
	"github.com/ninjadotorg/economy-simulation/state"
)

type Greedy struct {
	count int
}

func (g *Greedy) Init(a *agent.Agent) {
	log.Println("init greedy household agent")
}

func (g *Greedy) Run(a *agent.Agent, s state.State) {
	g.count++
	log.Println("i'm a greedy household", g.count)
	// rpc
}
