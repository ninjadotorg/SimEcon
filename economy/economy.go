package economy

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/ninjadotorg/economy-simulation/agent"
	"github.com/ninjadotorg/economy-simulation/agent/firm"
	"github.com/ninjadotorg/economy-simulation/agent/household"
	"github.com/ninjadotorg/economy-simulation/state"
)

const (
	INTERVAL = 1 * time.Second
)

type Economy struct {
	Groups []agent.Group `json:"agents"`
	agents []agent.Agent
}

func new(file string) (econ Economy, e error) {
	if f, e := ioutil.ReadFile(file); e != nil {
		return econ, e
	} else if e = json.Unmarshal(f, &econ); e != nil {
		return econ, e
	}

	for _, g := range econ.Groups {
		agent := agent.New(action(g.Action), g.StepSize)
		for i := 0; i < g.Qty; i++ {
			econ.agents = append(econ.agents, agent)
		}
	}
	return
}

func action(name string) (action agent.Action) {
	if name == "household.Greedy" {
		action = &household.Greedy{}
	} else if name == "household.Modest" {
		action = &household.Modest{}
	} else if name == "firm.Store" {
		action = &firm.Store{}
	}
	return action
}

func Run(file string) (e error) {
	if econ, e := new(file); e != nil {
		return e
	} else {

		// start all agents
		for i, _ := range econ.agents {
			go econ.agents[i].Run()
		}

		// broadcast state (loop)
		for step := 0; ; step++ {
			s := state.CurrentState()
			for _, a := range econ.agents {
				if step%a.StepSize == 0 {
					a.State <- s
				}
			}
			time.Sleep(INTERVAL)
		}
	}
}
