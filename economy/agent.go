package economy

import (
	"log"

	"github.com/ninjadotorg/SimEcon/util"
)

const (
	STATE_CHANNEL_SIZE    = 10
	CONTRACT_CHANNEL_SIZE = 10
	HOUR_CHANNEL_SIZE     = 10
)

type Group struct {
	Behavior string `json:"behavior"`
	Qty      int    `json:"qty"`
	StepSize int    `json:"stepSize"`
}

type Agent struct {
	behavior string
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

func newAgent(g Group) (a Agent) {
	a.uuid = util.NewUUID()
	a.stepSize = g.StepSize

	a.macro = make(chan State, STATE_CHANNEL_SIZE)
	a.contract = make(chan Contract, CONTRACT_CHANNEL_SIZE)
	a.hour = make(chan int, HOUR_CHANNEL_SIZE)

	a.action = action(g.Behavior)
	a.behavior = g.Behavior
	a.action.init(&a)

	return
}

func (a *Agent) run(econ *Economy) {
	for {
		select {

		case c := <-a.contract:
			log.Println("new contract", a.behavior, util.Shorten(a.uuid))
			a.action.handleContract(a, c, econ)

		case h := <-a.hour:
			// receive a clock reminder
			log.Println("new hourly checkup", a.behavior, util.Shorten(a.uuid))
			a.action.checkup(a, h, econ)

		case s := <-a.macro:
			// receive a (global) new network state update
			log.Println("new macro state", a.behavior, util.Shorten(a.uuid))
			a.action.run(a, s, econ)

		case <-a.quit:
			return

		}

	}
}
