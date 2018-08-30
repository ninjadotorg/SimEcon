package economy

import (
	"log"
)

type DefaultFirm struct {
	employees []Employee
}

type Employee struct {
	salary float64
	output float64
}

func (f *DefaultFirm) init(a *Agent) {
	log.Println("init")
}

func (f *DefaultFirm) run(a *Agent, s State, econ *Economy) {
}

func (f *DefaultFirm) handleContract(a *Agent, c Contract, econ *Economy) {
	if c.status == 1 && c.to.uuid == a.uuid {
		log.Println("accepting a new employee")
		c.status = 2
		a.handshake(c)
		econ.contracts = append(econ.contracts, c)
	}
}

func (f *DefaultFirm) checkup(a *Agent, hour int, econ *Economy) {
}
