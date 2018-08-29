package economy

import (
	"log"
)

type Greedy struct {
	count int
}

func (g *Greedy) init(a *Agent) {
	log.Println("init greedy household agent")
}

func (g *Greedy) run(a *Agent, s State) {
	g.count++
	log.Println("i'm a greedy household", g.count)
	// rpc
}
