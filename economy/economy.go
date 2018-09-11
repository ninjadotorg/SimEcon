package economy

import (
	"time"
)

const (
	TICK = time.Second
)

type AssetType string

type Economy struct {
	markets   map[AssetType]Market
	agents    []Agent
	contracts []Contract
}

func Run(file string) (e error) {

	var econ Economy
	if econ, e = newEconomy(file); e != nil {
		return e
	}

	// start all agents in separate goroutines
	for i, _ := range econ.agents {
		go econ.agents[i].run()
	}

	// clock tick
	for tick := 0; ; tick++ {
		for i, _ := range econ.agents {
			econ.agents[i].tick <- tick
		}
		time.Sleep(TICK)
	}

	// broadcast state (loop)
	// for step := 0; ; step++ {
	// 	s := currentState(econ)
	// 	for _, a := range econ.agents {
	// 		if step%a.stepSize == 0 {
	// 			a.macro <- s
	// 		}
	// 	}
	// 	time.Sleep(INTERVAL)
	// }

}
