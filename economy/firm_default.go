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
	log.Println("init factory")
}

func (f *DefaultFirm) run(a *Agent, s State, econ Economy) {
	log.Println("i'm a factory")
}

func (f *DefaultFirm) handleContract(a *Agent, c Contract, econ Economy) {
}

func (f *DefaultFirm) checkup(a *Agent, hour int, econ Economy) {

}
