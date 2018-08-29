package economy

import (
	"github.com/ninjadotorg/SimEcon/util"
)

const (
	CHANNEL_SIZE = 10
	PROBABILIY   = 80
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
	state    chan State
	contract chan Contract
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
	a.state = make(chan State, CHANNEL_SIZE)

	a.action = action
	a.action.init(&a)

	return
}

func (a *Agent) run() {
	for {
		select {

		case s := <-a.state:
			a.action.run(a, s)

			/*
				case c := <-a.contract:
					if rand.Intn(100) < PROBABILIY {

					}
			*/

		case <-a.quit:
			return

		}

	}
}
