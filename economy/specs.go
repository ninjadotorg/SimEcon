package economy

import (
	"encoding/json"
	"io/ioutil"
)

type EconomySpecs struct {
	Agents []AgentSpecs `json:"agents"`
}

type AgentSpecs struct {
	AgentType string `json:"agentType"`
	Qty       int    `json:"qty"`
}

func newEconomy(file string) (econ Economy, e error) {
	specs := EconomySpecs{}
	if f, e := ioutil.ReadFile(file); e != nil {
		return econ, e
	} else if e = json.Unmarshal(f, &specs); e != nil {
		return econ, e
	}

	for _, g := range specs.Agents {
		agent := newAgent(g, &econ)
		for i := 0; i < g.Qty; i++ {
			econ.agents = append(econ.agents, agent)
		}
	}
	return
}
