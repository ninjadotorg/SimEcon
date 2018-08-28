package economy

import (
	"encoding/json"
	"io/ioutil"
	"time"

	"github.com/0xroc/economy-simulation/agent"
	"github.com/0xroc/economy-simulation/state"
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
		agent := agent.New(g.Action, g.StepSize)
		for i := 0; i < g.Qty; i++ {
			econ.agents = append(econ.agents, agent)
		}
	}
	return
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
