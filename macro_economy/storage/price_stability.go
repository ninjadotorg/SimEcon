package storage

import (
	"fmt"

	"github.com/ninjadotorg/SimEcon/macro_economy/abstraction"
	"github.com/ninjadotorg/SimEcon/macro_economy/dto"
)

type PriceStability struct {
	Agent
}

func (ps *PriceStability) InitAgentAssets(
	st abstraction.Storage,
) {
	st.UpdateAssets(ps.AgentID, []abstraction.Asset{})
}

func (ps *PriceStability) GetType() uint {
	return ps.Type
}

func (ps *PriceStability) GetAgentID() string {
	return ps.AgentID
}

func (ps *PriceStability) UpdateAgent(
	st abstraction.Storage,
	agentDTO abstraction.AgentDTO,
) {
	updatingPriceStability := agentDTO.(dto.PriceStability)
	fmt.Println("updatingPriceStability: ", updatingPriceStability)
}
