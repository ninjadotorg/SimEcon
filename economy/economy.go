package economy

import (
	"encoding/json"
	"io/ioutil"
	"time"
)

const (
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
		agent := newAgent(action(g.Action), g.StepSize)
		for i := 0; i < g.Qty; i++ {
			econ.agents = append(econ.agents, agent)
		}
	}
	return
}

func action(name string) (action Action) {
	if name == "household.Greedy" {
		action = &Greedy{}
	} else if name == "household.Modest" {
		action = &Modest{}
	} else if name == "firm.Store" {
		action = &Restaurant{}
	}
	return action
}

func Run(file string) (e error) {
	if econ, e := newEconomy(file); e != nil {
		return e
	} else {

		// start all agents
		for i, _ := range econ.agents {
			go econ.agents[i].run()
		}

		// broadcast state (loop)
		for step := 0; ; step++ {
			s := currentState(econ)
			for _, a := range econ.agents {
				if step%a.stepSize == 0 {
					a.state <- s
				}
			}
			time.Sleep(INTERVAL)
		}
	}
}
