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

func (f *DefaultFirm) onContract(a *Agent, c Contract) {
	if c.status == 1 && c.to.uuid == a.uuid {
		log.Println("accepting a new employee")
		c.status = 2
		a.handshake(c)
		a.econ.contracts = append(a.econ.contracts, c)
	}
}

func (f *DefaultFirm) onTick(a *Agent, tick int) {
	m := a.econ.markets[LABOR_MARKET]
	a.buy(1, m.bestAsk())
}

func (f *DefaultFirm) produce(a *Agent) {

}
