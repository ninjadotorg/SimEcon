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

func (d *DefaultHousehold) onContract(a *Agent, c Contract) {
	if c.status == 2 && c.contractType == 1 {
		// employment contract
		log.Println("got a job")
		d.employed = 2
	}
}

func (d *DefaultHousehold) onTick(a *Agent, hour int) {
	if d.employed == 0 {

		for {

			firm := &a.econ.agents[rand.Intn(len(a.econ.agents))]

			if firm.agentType == "firm.default" {
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

func (d *DefaultHousehold) produce(a *Agent) {

}
