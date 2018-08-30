package economy

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

const (
	AN_HOUR  = 1 * time.Second // reduce 1h to 1s
	INTERVAL = 1 * time.Second
)

type Economy struct {
	Groups    []Group `json:"agents"`
	agents    []Agent
	contracts []Contract
}

func newEconomy(file string) (econ Economy, e error) {
	if f, e := ioutil.ReadFile(file); e != nil {
		return econ, e
	} else if e = json.Unmarshal(f, &econ); e != nil {
		return econ, e
	}

	for _, g := range econ.Groups {
		agent := newAgent(g)
		for i := 0; i < g.Qty; i++ {
			econ.agents = append(econ.agents, agent)
		}
	}
	return
}

func action(name string) (action Action) {
	if name == "household.default" {
		action = &DefaultHousehold{}
	} else if name == "firm.default" {
		action = &DefaultFirm{}
	} else if name == "firm.restaurant" {
		action = &Restaurant{}
	} else if name == "network.contract" {
		action = &ContractExecutor{}
	}
	return
}

func Run(file string) (e error) {

	var econ Economy
	if econ, e = newEconomy(file); e != nil {
		return e
	}

	// start all agents in separate goroutines
	for i, _ := range econ.agents {
		go econ.agents[i].run(&econ)
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

	// clock broadcast hourly
	for hour := 0; ; hour++ {
		for i, _ := range econ.agents {
			econ.agents[i].hour <- hour
		}
		time.Sleep(AN_HOUR)
	}

	// order book approach

}
