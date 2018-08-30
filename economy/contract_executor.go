package economy

import (
	"log"

	"github.com/ninjadotorg/SimEcon/util"
)

type ContractExecutor struct {
	employed int // 0: unemployed, 1: applying, 2: employed
	salary   float64
}

func (d *ContractExecutor) init(a *Agent) {
}

func (d *ContractExecutor) run(a *Agent, s State, econ *Economy) {
}

func (d *ContractExecutor) handleContract(a *Agent, c Contract, econ *Economy) {
}

func (d *ContractExecutor) checkup(a *Agent, hour int, econ *Economy) {

	// executing contract
	for _, c := range econ.contracts {
		if hour%c.repeat == 0 {
			c.payer.asset.cash -= c.amt
			c.payee.asset.cash += c.amt
		}
	}

	// print agent balances
	for _, a := range econ.agents {
		log.Println(util.Shorten(a.uuid), a.behavior, a.asset.cash)
	}
}
