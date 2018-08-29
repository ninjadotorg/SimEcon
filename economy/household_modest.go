package economy

import (
	"log"
)

type Modest struct {
	// state
}

func (m *Modest) init(a *Agent) {
	log.Println("init modest household agent")
}

func (m *Modest) run(a *Agent, s State) {
	log.Println("i'm a modest household")
}
