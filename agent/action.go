package agent

import "github.com/ninjadotorg/economy-simulation/state"

type Action interface {
	Init(a *Agent)
	Run(a *Agent, s state.State)
}
