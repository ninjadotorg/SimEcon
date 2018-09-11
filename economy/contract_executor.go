package economy

import (
	"log"

	"github.com/ninjadotorg/SimEcon/util"
)

type ContractExecutor struct {
	employed int // 0: unemployed, 1: applying, 2: employed
	salary   float64
}

func (d *ContractExecutor) onContract(a *Agent, c Contract) {
}

func (d *ContractExecutor) onTick(a *Agent, tick int) {

	// executing contract
	for _, c := range a.econ.contracts {
		if tick%c.repeat == 0 {
			c.payer.cash -= c.amt
			c.payee.cash += c.amt
		}
	}

	// print agent balances
	for _, a := range a.econ.agents {
		log.Println(util.Shorten(a.uuid), a.agentType, a.cash)
	}
}

func (d *ContractExecutor) produce(a *Agent) {

}
