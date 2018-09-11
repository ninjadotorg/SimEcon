package economy

import (
	"errors"
	"log"

	"github.com/ninjadotorg/SimEcon/util"
)

const (
	STATE_CHANNEL_SIZE    = 10
	CONTRACT_CHANNEL_SIZE = 10
	HOUR_CHANNEL_SIZE     = 10
)

type Agent struct {
	agentType string
	policy    Policy
	uuid      string
	buyAsset  Asset
	sellAsset Asset
	cash      float64
	econ      *Economy

	Communications
}

type Asset struct {
	assetType   AssetType
	targetPrice float64
	inventory   int
}

type Communications struct {
	tick     chan int      // the smallest unit of time in this economy
	contract chan Contract // p2p direct contract w/ another agent
	macro    chan State    // network state broadcased to all agents by the network
	quit     chan struct{}
}

// type BalanceSheet struct {

// 	// assets
// 	quantity []float64
// 	price    []float64
// 	cash     float64

// 	// liabilities
// 	equity float64
// 	debt   float64
// }

func newAgent(s AgentSpecs, econ *Economy) (a Agent) {
	a.uuid = util.NewUUID()

	a.macro = make(chan State, STATE_CHANNEL_SIZE)
	a.contract = make(chan Contract, CONTRACT_CHANNEL_SIZE)
	a.tick = make(chan int, HOUR_CHANNEL_SIZE)

	a.agentType = s.AgentType
	a.policy = policy(s.AgentType)
	a.econ = econ

	return
}

func (a *Agent) run() {
	for {
		select {

		case tick := <-a.tick:
			log.Println("new tick", a.agentType, util.Shorten(a.uuid))
			a.policy.onTick(a, tick)

		// case contract := <-agent.contract:
		// 	log.Println("new contract", agent.agentType, util.Shorten(agent.uuid))
		// 	agent.policy.onContract(agent, contract, econ)

		// case macro := <-a.macro:
		// 	log.Println("new macro state", a.agentType, util.Shorten(a.uuid))
		// 	a.policy.run(agent, macro, econ)

		case <-a.quit:
			return

		}

	}
}

func (a *Agent) buy(size float64, price float64) error {
	if size*price > a.cash {
		return errors.New("not enough cash")
	}
	m := a.econ.markets[a.buyAsset.assetType]
	m.trade("buy", size, price, a.uuid)
	return nil
}

func (a *Agent) sell(size float64, price float64) {
	m := a.econ.markets[a.sellAsset.assetType]
	m.trade("sell", size, price, a.uuid)
}

func (a *Agent) produce() {
	a.policy.produce(a)
}
