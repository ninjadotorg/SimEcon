package market

import (
	"fmt"
	"math"
	"time"

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
	prod abstraction.Production,
	tr abstraction.Tracker,
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
		sellerActualAsset := prod.GetActualAsset(sellerAsset)
		actualBidQty := math.Min(sellerActualAsset.GetQuantity(), bid.GetQuantity())
		if actualBidQty >= orderItemReq.Quantity {
			am.Pay(
				agentID,
				bid.GetAgentID(),
				bid.GetPricePerUnit()*orderItemReq.Quantity,
				common.PRIIC,
			)
			bid.SetQuantity(actualBidQty - orderItemReq.Quantity)
			if bid.GetQuantity() == 0 {
				removingBidAgentIDs = append(removingBidAgentIDs, bid.GetAgentID())
			}
			sellerActualAsset.SetQuantity(sellerActualAsset.GetQuantity() - orderItemReq.Quantity)
			st.UpdateAsset(bid.GetAgentID(), sellerActualAsset)
			orderItemReq.Quantity = 0
			break
		}
		am.Pay(
			agentID,
			bid.GetAgentID(),
			bid.GetPricePerUnit()*actualBidQty,
			common.PRIIC,
		)
		orderItemReq.Quantity -= actualBidQty
		sellerActualAsset.SetQuantity(sellerActualAsset.GetQuantity() - actualBidQty)
		st.UpdateAsset(bid.GetAgentID(), sellerActualAsset)

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
		totalAsks := st.GetTotalAsksByAssetType(orderItemReq.AssetType)
		record := []string{fmt.Sprintf("%d", time.Now().Unix()), fmt.Sprintf("%.1f", totalAsks)}

		err := tr.WriteToCSV(fmt.Sprintf("%s_%d.csv", common.TOTAL_ASKS_FILE, orderItemReq.AssetType), record)
		if err != nil {
			return orderItemReq.Quantity, err
		}
	}

	return orderItemReq.Quantity, nil
}

func (m *Market) Sell(
	agentID string,
	orderItemReq *dto.OrderItem,
	st abstraction.Storage,
	am abstraction.AccountManager,
	prod abstraction.Production,
	tr abstraction.Tracker,
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
		buyerActualAsset := prod.GetActualAsset(buyerAsset)

		askBalance := am.GetBalance(ask.GetAgentID())
		if askBalance < ask.GetPricePerUnit()*math.Min(orderItemReq.Quantity, ask.GetQuantity()) {
			removingAskAgentIDs = append(removingAskAgentIDs, ask.GetAgentID())
			continue
		}

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

			buyerActualAsset.SetQuantity(buyerActualAsset.GetQuantity() + orderItemReq.Quantity)
			st.UpdateAsset(ask.GetAgentID(), buyerActualAsset)

			orderItemReq.Quantity = 0
			break
		}
		am.Pay(
			ask.GetAgentID(),
			agentID,
			ask.GetPricePerUnit()*ask.GetQuantity(),
			common.PRIIC,
		)

		buyerActualAsset.SetQuantity(buyerActualAsset.GetQuantity() + ask.GetQuantity())
		st.UpdateAsset(ask.GetAgentID(), buyerActualAsset)

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
		totalBids := st.GetTotalBidsByAssetType(orderItemReq.AssetType)
		record := []string{fmt.Sprintf("%d", time.Now().Unix()), fmt.Sprintf("%.1f", totalBids)}
		fmt.Println("Record: ", record)
		err := tr.WriteToCSV(fmt.Sprintf("%s_%d.csv", common.TOTAL_BIDS_FILE, orderItemReq.AssetType), record)
		if err != nil {
			return orderItemReq.Quantity, err
		}
	}

	return orderItemReq.Quantity, nil
}
