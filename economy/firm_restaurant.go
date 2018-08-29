package economy

import (
	"log"
)

type Restaurant struct{}

func (r *Restaurant) init(a *Agent) {
	log.Println("init restaurant")
}

func (r *Restaurant) run(a *Agent, s State) {
	log.Println("i'm a restaurant")
}
