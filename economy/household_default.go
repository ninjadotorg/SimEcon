package economy

import (
	"log"
	"math/rand"
)

const (
	MONTHLY = 1 * 24 * 30
)

type DefaultHousehold struct {
	employed int // 0: unemployed, 1: applying, 2: employed
	salary   float64
}

func (d *DefaultHousehold) init(a *Agent) {
	log.Println("init default household agent")
}

func (d *DefaultHousehold) run(a *Agent, s State, econ Economy) {
}

func (d *DefaultHousehold) handleContract(a *Agent, c Contract, econ Economy) {
	if c.status == 1 && c.to == a {
		c.status = 2
		a.handshake(c)
		econ.contracts = append(econ.contracts, c)
	} else if c.status == 2 && c.contractType == 1 {
		// employment contract
		d.employed = 2
	}
}

func (d *DefaultHousehold) checkup(a *Agent, hour int, econ Economy) {
	log.Println("i'm a default household agent")
	if d.employed == 0 {
		// look for a job
		for {
			firm := econ.agents[rand.Intn(len(econ.agents))]
			if firm.agentType == 1 {
				d.employed = 1
				c := Contract{}
				c.from = a
				c.to = &firm
				c.amt = d.salary
				c.repeat = MONTHLY
				a.initiate(c)
				break
			}
		}
	}
}
