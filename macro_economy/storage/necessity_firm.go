package storage

import (
	"time"

	"github.com/ninjadotorg/SimEcon/common"
	"github.com/ninjadotorg/SimEcon/macro_economy/abstraction"
)

type NecessityFirm struct {
	Agent
}

func (nf *NecessityFirm) InitAgentAssets(
	st abstraction.Storage,
) {
	// man hours asset
	mhAsset := &Asset{
		AgentID:      nf.AgentID,
		Type:         common.MAN_HOUR,
		Quantity:     common.NFIRM_MAN_HOURS,
		ProducedTime: time.Now().Unix(),
	}
	// capital asset
	cAsset := &Asset{
		AgentID:      nf.AgentID,
		Type:         common.CAPITAL,
		Quantity:     common.NFIRM_CAPITAL,
		ProducedTime: time.Now().Unix(),
	}
	// necessity asset
	nAsset := &Asset{
		AgentID:      nf.AgentID,
		Type:         common.NECESSITY,
		Quantity:     common.NFIRM_NECESSITY,
		ProducedTime: time.Now().Unix(),
	}
	st.UpdateAssets(nf.AgentID, []abstraction.Asset{mhAsset, cAsset, nAsset})
}

func (nf *NecessityFirm) GetType() uint {
	return nf.Type
}
