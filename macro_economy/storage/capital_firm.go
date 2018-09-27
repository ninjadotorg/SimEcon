package storage

import (
	"fmt"
	"time"

	"github.com/ninjadotorg/SimEcon/common"
	"github.com/ninjadotorg/SimEcon/macro_economy/abstraction"
	"github.com/ninjadotorg/SimEcon/macro_economy/dto"
)

type CapitalFirm struct {
	Agent
}

func (c *CapitalFirm) InitAgentAssets(
	st abstraction.Storage,
) {
	// capital asset
	cAsset := &Asset{
		AgentID:      c.AgentID,
		Type:         common.CAPITAL,
		Quantity:     common.CAPITAL_CAPITAL,
		ProducedTime: time.Now().Unix(),
	}
	// man hours asset
	mhAsset := &Asset{
		AgentID:      c.AgentID,
		Type:         common.MAN_HOUR,
		Quantity:     common.CAPITAL_MAN_HOURS,
		ProducedTime: time.Now().Unix(),
	}
	st.UpdateAssets(c.AgentID, []abstraction.Asset{mhAsset, cAsset})
}

func (c *CapitalFirm) GetType() uint {
	return c.Type
}

func (c *CapitalFirm) UpdateAgent(
	st abstraction.Storage,
	agentDTO abstraction.AgentDTO,
) {
	updatingCFirm := agentDTO.(*dto.CapitalFirm)
	fmt.Println("updatingCFirm: ", updatingCFirm)
}
