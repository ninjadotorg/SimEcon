package storage

import (
	"fmt"
	"time"

	"github.com/ninjadotorg/SimEcon/common"
	"github.com/ninjadotorg/SimEcon/macro_economy/abstraction"
	"github.com/ninjadotorg/SimEcon/macro_economy/dto"
)

type Miner struct {
	Agent
	// ServerID    string
	// MinerStatus uint
	// BlockHeight uint
}

func (m *Miner) InitAgentAssets(
	st abstraction.Storage,
) {
	// necessity asset
	nAsset := &Asset{
		AgentID:      m.AgentID,
		Type:         common.NECESSITY,
		Quantity:     common.MINER_NECESSITY,
		ProducedTime: time.Now().Unix(),
	}

	st.UpdateAssets(m.AgentID, []abstraction.Asset{nAsset})
}

func (m *Miner) GetType() uint {
	return m.Type
}

func (m *Miner) GetAgentID() string {
	return m.AgentID
}

func (m *Miner) UpdateAgent(
	st abstraction.Storage,
	agentDTO abstraction.AgentDTO,
) {
	updatingMiner := agentDTO.(dto.Miner)
	fmt.Println("updatingMiner: ", updatingMiner)
}
