package economy

import (
	"log"
)

type Restaurant struct{}

func (r *Restaurant) init(a *Agent) {
	log.Println("init")
}

func (r *Restaurant) run(a *Agent, s State, econ *Economy) {
}

func (r *Restaurant) handleContract(a *Agent, c Contract, econ *Economy) {
	log.Println("HERE")
}

func (r *Restaurant) checkup(a *Agent, hour int, econ *Economy) {
}
