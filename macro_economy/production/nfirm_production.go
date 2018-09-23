package production

import (
	"github.com/ninjadotorg/SimEcon/common"
	"github.com/ninjadotorg/SimEcon/macro_economy/abstraction"
	"github.com/ninjadotorg/SimEcon/macro_economy/dto"
)

type NFirmProduction struct {
}

func (nfp *NFirmProduction) Produce(
	st abstraction.Storage,
	agentID string,
	assetsReq map[uint]*dto.Asset,
) (map[uint]abstraction.Asset, error) {
	agentAssets, _ := st.GetAgentAssets(agentID)
	nAsset := agentAssets[common.NECESSITY]
	mhAsset := agentAssets[common.MAN_HOUR]
	cAsset := agentAssets[common.CAPITAL]

	cAssetReq := assetsReq[common.CAPITAL]
	mhAssetReq := assetsReq[common.MAN_HOUR]

	curNAsset := computeDecayNecessity(nAsset)
	curMHAsset := computeDecayManHours(mhAsset)
	curCAsset := computeDecayCapital(cAsset)

	spendingCAmt := cAssetReq.Quantity
	if spendingCAmt > curCAsset.GetQuantity() {
		spendingCAmt = curCAsset.GetQuantity()
	}
	spendingMHAmt := mhAssetReq.Quantity
	if spendingMHAmt > curMHAsset.GetQuantity() {
		spendingMHAmt = curMHAsset.GetQuantity()
	}

	// TODO: define convertion formula
	input := (spendingMHAmt + spendingCAmt) / 2
	convertedNAmt := convertLinearly(input, 1.5)

	curNAsset.SetQuantity(convertedNAmt + curNAsset.GetQuantity())
	curNAsset.SetProducedTime()

	curMHAsset.SetQuantity(curMHAsset.GetQuantity() - spendingMHAmt)
	curMHAsset.SetProducedTime()

	curCAsset.SetQuantity(curCAsset.GetQuantity() - spendingCAmt)
	curCAsset.SetProducedTime()

	st.UpdateAsset(agentID, curNAsset)
	st.UpdateAsset(agentID, curMHAsset)
	st.UpdateAsset(agentID, curCAsset)

	res, _ := st.GetAgentAssets(agentID)
	return res, nil
}

func (nfp *NFirmProduction) GetActualAsset(
	asset abstraction.Asset,
) abstraction.Asset {
	return computeDecayNecessity(asset)
}
