package agent

import (
	"github.com/0xroc/economy-simulation/agent/firm"
	"github.com/0xroc/economy-simulation/agent/household"
	"github.com/0xroc/economy-simulation/state"
	"github.com/0xroc/economy-simulation/util"
)

type Agent struct {
	StepSize int
	State    chan state.State
	Quit     chan struct{}

	uuid string
	act  func(state.State)
}

type Group struct {
	Action   string `json:"action"`
	Qty      int    `json:"qty"`
	StepSize int    `json:"stepSize"`
}

func New(actionName string, stepSize int) (a Agent) {
	a.uuid = util.NewUUID()
	a.StepSize = stepSize
	a.act = action(actionName)
	a.State = make(chan state.State, 10)
	return
}

func action(name string) func(state.State) {
	if name == "household.Greedy" {
		return household.Greedy
	} else if name == "household.Modest" {
		return household.Modest
	} else if name == "firm.Store" {
		return firm.Store
	}
	return nil
}

func (a *Agent) Run() {
	for {
		a.act(<-a.State)
	}
}
