package agent

import (
	"github.com/ninjadotorg/economy-simulation/state"
	"github.com/ninjadotorg/economy-simulation/util"
)

const (
	CHANNEL_SIZE = 10
)

type Agent struct {
	StepSize  int
	State     chan state.State
	Quit      chan struct{}
	action    Action
	uuid      string
	asset     Asset
	liability Liability
	state     int // state of an agent: running, paused
}

type Asset struct {
	quantity []float64
	price    []float64
	cash     float64
}

type Liability struct {
	equity float64
	debt   float64
}

type Group struct {
	Action   string `json:"action"`
	Qty      int    `json:"qty"`
	StepSize int    `json:"stepSize"`
}

func New(action Action, stepSize int) (a Agent) {
	a.uuid = util.NewUUID()
	a.StepSize = stepSize
	a.State = make(chan state.State, CHANNEL_SIZE)

	a.action = action
	a.action.Init(&a)

	return
}

func (a *Agent) Run() {
	for {
		select {
		case state := <-a.State:
			a.action.Run(a, state)
		case <-a.Quit:
			return
		}
	}
}
