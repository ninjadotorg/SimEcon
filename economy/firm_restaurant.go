package economy

import (
	"log"
)

type Restaurant struct{}

func (r *Restaurant) init(a *Agent) {
	log.Println("init restaurant")
}

func (r *Restaurant) run(a *Agent, s State, econ Economy) {
}

func (r *Restaurant) handleContract(a *Agent, c Contract, econ Economy) {
}

func (r *Restaurant) checkup(a *Agent, hour int, econ Economy) {
	log.Println("i'm a restaurant")

}
