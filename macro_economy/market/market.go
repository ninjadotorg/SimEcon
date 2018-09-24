package market

import (
	"github.com/ninjadotorg/SimEcon/common"
	"github.com/ninjadotorg/SimEcon/macro_economy/abstraction"
	"github.com/ninjadotorg/SimEcon/macro_economy/dto"
)

type Market struct {
}

var market *Market

func GetMarketInstance() *Market {
	if market != nil {
		return market
	}
	market = &Market{}
	return market
}

func (m *Market) Buy(
	agentID string,
	orderItemReq *dto.OrderItem,
	st abstraction.Storage,
	am abstraction.AccountManager,
) (float64, error) {
	sortedBidsByAssetType := st.GetSortedBidsByAssetType(orderItemReq.AssetType, false)

	removingBidAgentIDs := []string{}
	for _, bid := range sortedBidsByAssetType {
		if bid.GetPricePerUnit() > orderItemReq.PricePerUnit {
			continue
		}

		sellerAsset, _ := st.GetAgentAsset(
			bid.GetAgentID(),
			bid.GetAssetType(),
		)
		// TODO: get actual qty of asset of agent who owns the bid
		if bid.GetQuantity() >= orderItemReq.Quantity {
			am.Pay(
				agentID,
				bid.GetAgentID(),
				bid.GetPricePerUnit()*orderItemReq.Quantity,
				common.PRIIC,
			)
			bid.SetQuantity(bid.GetQuantity() - orderItemReq.Quantity)
			if bid.GetQuantity() == 0 {
				removingBidAgentIDs = append(removingBidAgentIDs, bid.GetAgentID())
			}
			sellerAsset.SetQuantity(sellerAsset.GetQuantity() - orderItemReq.Quantity)
			st.UpdateAsset(bid.GetAgentID(), sellerAsset)
			orderItemReq.Quantity = 0
			break
		}
		am.Pay(
			agentID,
			bid.GetAgentID(),
			bid.GetPricePerUnit()*bid.GetQuantity(),
			common.PRIIC,
		)
		sellerAsset.SetQuantity(sellerAsset.GetQuantity() - bid.GetQuantity())
		st.UpdateAsset(bid.GetAgentID(), sellerAsset)

		orderItemReq.Quantity -= bid.GetQuantity()
		bid.SetQuantity(0)
		removingBidAgentIDs = append(removingBidAgentIDs, bid.GetAgentID())
	}
	// re-update bid list: remove bid with qty = 0 and append new ask if remaning qty > 0
	if len(removingBidAgentIDs) > 0 {
		err := st.RemoveBidsByAgentIDs(removingBidAgentIDs, orderItemReq.AssetType)
		if err != nil {
			return -1, err
		}
	}

	if orderItemReq.Quantity > 0 {
		st.AppendAsk(
			orderItemReq.AssetType,
			orderItemReq.AgentID,
			orderItemReq.Quantity,
			orderItemReq.PricePerUnit,
		)
	}

	return orderItemReq.Quantity, nil
}

func (m *Market) Sell(
	agentID string,
	orderItemReq *dto.OrderItem,
	st abstraction.Storage,
	am abstraction.AccountManager,
) (float64, error) {
	sortedAsksByAssetType := st.GetSortedAsksByAssetType(orderItemReq.AssetType, false)

	removingAskAgentIDs := []string{}
	for _, ask := range sortedAsksByAssetType {
		if ask.GetPricePerUnit() < orderItemReq.PricePerUnit {
			continue
		}

		buyerAsset, _ := st.GetAgentAsset(
			ask.GetAgentID(),
			ask.GetAssetType(),
		)

		if ask.GetQuantity() >= orderItemReq.Quantity {
			am.Pay(
				ask.GetAgentID(),
				agentID,
				ask.GetPricePerUnit()*orderItemReq.Quantity,
				common.PRIIC,
			)
			ask.SetQuantity(ask.GetQuantity() - orderItemReq.Quantity)
			if ask.GetQuantity() == 0 {
				removingAskAgentIDs = append(removingAskAgentIDs, ask.GetAgentID())
			}

			buyerAsset.SetQuantity(buyerAsset.GetQuantity() + orderItemReq.Quantity)
			st.UpdateAsset(ask.GetAgentID(), buyerAsset)

			orderItemReq.Quantity = 0
			break
		}
		am.Pay(
			ask.GetAgentID(),
			agentID,
			ask.GetPricePerUnit()*ask.GetQuantity(),
			common.PRIIC,
		)

		buyerAsset.SetQuantity(buyerAsset.GetQuantity() + ask.GetQuantity())
		st.UpdateAsset(ask.GetAgentID(), buyerAsset)

		orderItemReq.Quantity -= ask.GetQuantity()
		ask.SetQuantity(0)
		removingAskAgentIDs = append(removingAskAgentIDs, ask.GetAgentID())
	}
	// re-update ask list: remove ask with qty = 0 and append new ask if remaning qty > 0
	if len(removingAskAgentIDs) > 0 {
		err := st.RemoveAsksByAgentIDs(removingAskAgentIDs, orderItemReq.AssetType)
		if err != nil {
			return -1, err
		}
	}

	if orderItemReq.Quantity > 0 {
		st.AppendBid(
			orderItemReq.AssetType,
			orderItemReq.AgentID,
			orderItemReq.Quantity,
			orderItemReq.PricePerUnit,
		)
	}

	return orderItemReq.Quantity, nil
}
