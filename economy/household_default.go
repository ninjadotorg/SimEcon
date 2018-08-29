package economy

import (
	"log"
)

type Default struct {
	employed bool
}

func (d *Default) init(a *Agent) {
	log.Println("init default household agent")
}

func (d *Default) run(a *Agent, s State) {
	log.Println("i'm a default household agent")
	if !d.employed {
		// look for a job
		for _, firm := range s.agents {
			if firm.agentType == 1 {
				a.initiate(firm, Contract{})
			}
		}
	}
}
