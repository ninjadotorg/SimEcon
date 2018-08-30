package economy

import (
	"log"
	"math/rand"

	"github.com/ninjadotorg/SimEcon/util"
)

const (
	MONTHLY = 1 * 24 * 30
)

type DefaultHousehold struct {
	employed int // 0: unemployed, 1: applying, 2: employed
	salary   float64
}

func (d *DefaultHousehold) init(a *Agent) {
	log.Println("init")
	d.salary = 1
}

func (d *DefaultHousehold) run(a *Agent, s State, econ *Economy) {
}

func (d *DefaultHousehold) handleContract(a *Agent, c Contract, econ *Economy) {
	if c.status == 2 && c.contractType == 1 {
		// employment contract
		log.Println("got a job")
		d.employed = 2
	}
}

func (d *DefaultHousehold) checkup(a *Agent, hour int, econ *Economy) {
	if d.employed == 0 {

		for {

			firm := &econ.agents[rand.Intn(len(econ.agents))]

			if firm.behavior == "firm.default" {
				log.Println("looking for a job at", util.Shorten(firm.uuid))
				d.employed = 1

				c := Contract{}
				c.from = a
				c.to = firm
				c.payer = firm
				c.payee = a

				c.contractType = 1 // employment

				c.status = 1
				c.amt = d.salary
				c.repeat = 1

				a.initiate(c)

				break
			}
		}
	}
}
