package production

import (
	"github.com/ninjadotorg/SimEcon/common"
	"github.com/ninjadotorg/SimEcon/macro_economy/abstraction"
	"github.com/ninjadotorg/SimEcon/macro_economy/dto"
)

type CFirmProduction struct {
}

func (cfp *CFirmProduction) Produce(
	st abstraction.Storage,
	agentID string,
	assetsReq map[uint]*dto.Asset,
) (map[uint]abstraction.Asset, error) {
	agentAssets, _ := st.GetAgentAssets(agentID)
	cAsset := agentAssets[common.CAPITAL]
	mhAsset := agentAssets[common.MAN_HOUR]

	mhAssetReq := assetsReq[common.MAN_HOUR]

	curCAsset := computeDecayCapital(cAsset)
	curMHAsset := computeDecayManHours(mhAsset)

	spendingAmt := mhAssetReq.Quantity
	if spendingAmt > curMHAsset.GetQuantity() {
		spendingAmt = curMHAsset.GetQuantity()
	}
	convertedCAmt := convertLinearly(spendingAmt, 0.75) // TODO
	curCAsset.SetQuantity(convertedCAmt + curCAsset.GetQuantity())
	curCAsset.SetProducedTime()

	curMHAsset.SetQuantity(curMHAsset.GetQuantity() - spendingAmt)
	curMHAsset.SetProducedTime()

	st.UpdateAsset(agentID, curMHAsset)
	st.UpdateAsset(agentID, curCAsset)

	res, _ := st.GetAgentAssets(agentID)
	return res, nil
}
