package sim

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"time"
)

const (
	INTERVAL = 1 * time.Second
)

type Sim struct {
	AgentGroups []AgentGroup `json:"agents"`
	agents      []Agent
}

func newSim(file string) (s Sim, e error) {
	if f, e := ioutil.ReadFile(file); e != nil {
		return s, e
	} else if e = json.Unmarshal(f, &s); e != nil {
		return s, e
	}

	for _, g := range s.AgentGroups {
		agent := newAgent(g.Type, g.StepSize)
		for i := 0; i < g.Qty; i++ {
			s.agents = append(s.agents, agent)
		}
	}
	return
}

func Start(file string) (e error) {
	if sim, e := newSim(file); e != nil {
		return e
	} else {

		log.Println("simulation:", sim)

		// start all agents
		for i, _ := range sim.agents {
			go sim.agents[i].start()
		}

		for step := 0; ; step++ {

			s := state()

			for _, a := range sim.agents {
				if step%a.stepSize == 0 {
					a.state <- s
				}
			}

			time.Sleep(INTERVAL)
		}
	}
}
