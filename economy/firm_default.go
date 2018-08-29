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

func (f *DefaultFirm) run(a *Agent, s State) {
	log.Println("i'm a factory")
}
