package economy

import (
	"github.com/ninjadotorg/SimEcon/util"
)

const (
	STATE_CHANNEL_SIZE    = 10
	CONTRACT_CHANNEL_SIZE = 10
	HOUR_CHANNEL_SIZE     = 10
)

type Agent struct {

	// 0: household (default)
	// 1: firm
	// 2: bank
	// 3: central bank
	agentType int

	stepSize int
	action   Action
	uuid     string

	// balance sheet
	asset     Asset
	liability Liability

	// communications
	contract chan Contract
	macro    chan State
	hour     chan int
	quit     chan struct{}
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

func newAgent(action Action, stepSize int) (a Agent) {
	a.uuid = util.NewUUID()
	a.stepSize = stepSize
	a.macro = make(chan State, STATE_CHANNEL_SIZE)
	a.contract = make(chan Contract, CONTRACT_CHANNEL_SIZE)
	a.hour = make(chan int, HOUR_CHANNEL_SIZE)

	a.action = action
	a.action.init(&a)

	return
}

func (a *Agent) run(econ Economy) {
	for {
		select {

		case h := <-a.hour:
			// receive a (global) new network state update
			a.action.checkup(a, h, econ)

		case s := <-a.macro:
			// receive a (global) new network state update
			a.action.run(a, s, econ)

		case c := <-a.contract:
			a.action.handleContract(a, c, econ)

		case <-a.quit:
			return

		}

	}
}
