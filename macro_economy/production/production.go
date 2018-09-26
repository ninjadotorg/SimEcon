package production

import (
	"errors"
	"math"
	"time"

	"github.com/ninjadotorg/SimEcon/common"
	"github.com/ninjadotorg/SimEcon/macro_economy/abstraction"
)

type Production struct {
	AgentTypeToAgentProduction map[uint]abstraction.AgentProduction
	AssetTypeToDecayFunction   map[uint]func(abstraction.Asset) abstraction.Asset
}

var production *Production

func GetProductionInstance() *Production {
	if production != nil {
		return production
	}
	production = &Production{
		AgentTypeToAgentProduction: map[uint]abstraction.AgentProduction{
			common.PERSON:         &PersonProduction{},
			common.NECESSITY_FIRM: &NFirmProduction{},
			common.CAPITAL_FIRM:   &CFirmProduction{},
		},
		AssetTypeToDecayFunction: map[uint]func(abstraction.Asset) abstraction.Asset{
			common.NECESSITY: computeDecayNecessity,
			common.CAPITAL:   computeDecayCapital,
			common.MAN_HOUR:  computeDecayManHours,
		},
	}
	return production
}

func computeDecayNecessity(asset abstraction.Asset) abstraction.Asset {
	decaySteps := int(math.Floor(float64(time.Now().Unix()-asset.GetProducedTime()) / common.NECESSITY_DECAY_PERIOD)) // decay period = 5m
	for i := 1; i <= decaySteps; i++ {
		asset.SetQuantity(common.NECESSITY_EPSILON_DECAY * asset.GetQuantity())
	}
	return asset
}

func computeDecayCapital(asset abstraction.Asset) abstraction.Asset {
	decaySteps := int(math.Floor(float64(time.Now().Unix()-asset.GetProducedTime()) / common.CAPITAL_DECAY_PERIOD)) // decay period = 4m
	for i := 1; i <= decaySteps; i++ {
		asset.SetQuantity(common.CAPITAL_EPSILON_DECAY * asset.GetQuantity())
	}
	return asset
}

func computeDecayManHours(asset abstraction.Asset) abstraction.Asset {
	decaySteps := int(math.Floor(float64(time.Now().Unix()-asset.GetProducedTime()) / common.MAN_HOURS_DECAY_PERIOD)) // decay period = 6m
	for i := 1; i <= decaySteps; i++ {
		asset.SetQuantity(common.MAN_HOURS_EPSILON_DECAY * asset.GetQuantity())
	}
	return asset
}

func convertLinearly(input float64, a float64) float64 {
	return a * input
}

func (prod *Production) GetProductionByAgentType(
	agentType uint,
) (abstraction.AgentProduction, error) {
	agentProd, ok := prod.AgentTypeToAgentProduction[agentType]
	if !ok {
		return nil, errors.New("Agent ID not found")
	}
	return agentProd, nil
}

func (prod *Production) GetActualAsset(
	asset abstraction.Asset,
) abstraction.Asset {
	assetType := asset.GetType()
	decayFunc := prod.AssetTypeToDecayFunction[assetType]
	return decayFunc(asset)
}

func (prod *Production) GetActualAssets(
	assets map[uint]abstraction.Asset,
) map[uint]abstraction.Asset {
	actualAssets := map[uint]abstraction.Asset{}
	for assetType, asset := range assets {
		decayFunc := prod.AssetTypeToDecayFunction[assetType]
		actualAssets[assetType] = decayFunc(asset)
	}
	return actualAssets
}
