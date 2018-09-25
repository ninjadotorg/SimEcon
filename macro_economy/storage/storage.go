package storage

import (
	"errors"
	"time"

	"github.com/ninjadotorg/SimEcon/common"
	"github.com/ninjadotorg/SimEcon/macro_economy/abstraction"
)

type Storage struct {
	Agents map[string]abstraction.Agent
	Assets map[string]map[uint]abstraction.Asset     // agentID -> assetID -> asset
	Asks   map[uint]map[string]abstraction.OrderItem // assetType -> agentID -> orderItem
	Bids   map[uint]map[string]abstraction.OrderItem
}

var storage *Storage

func GetStorageInstance() *Storage {
	if storage != nil {
		return storage
	}
	storage = &Storage{
		Agents: map[string]abstraction.Agent{},
		Assets: map[string]map[uint]abstraction.Asset{},
		Asks:   map[uint]map[string]abstraction.OrderItem{},
		Bids:   map[uint]map[string]abstraction.OrderItem{},
	}
	return storage
}

func (st *Storage) InsertAgent(
	agentID string,
	agentType uint,
) abstraction.Agent {
	agent := Agent{
		AgentID: agentID,
		Type:    agentType,
	}
	var newAgent abstraction.Agent = nil
	if agentType == common.PERSON {
		newAgent = &Person{
			agent,
		}
	} else if agentType == common.NECESSITY_FIRM {
		newAgent = &NecessityFirm{
			agent,
		}
	} else if agentType == common.CAPITAL_FIRM {
		newAgent = &CapitalFirm{
			agent,
		}
	}

	st.Agents[agentID] = newAgent
	return newAgent
}

func (st *Storage) GetAgentByID(
	agentID string,
) (abstraction.Agent, error) {

	agent, ok := st.Agents[agentID]
	if !ok {
		return nil, errors.New("Could not find the agent")
	}
	return agent, nil
}

func (st *Storage) UpdateAssets(
	agentID string,
	assets []abstraction.Asset,
) {
	if _, ok := st.Assets[agentID]; !ok {
		st.Assets[agentID] = map[uint]abstraction.Asset{}
	}
	for _, asset := range assets {
		st.Assets[agentID][asset.GetType()] = asset
	}
}

func (st *Storage) UpdateAsset(
	agentID string,
	asset abstraction.Asset,
) {
	if _, ok := st.Assets[agentID]; !ok {
		st.Assets[agentID] = map[uint]abstraction.Asset{}
	}
	st.Assets[agentID][asset.GetType()] = asset
}

func (st *Storage) GetAgentAssets(
	agentID string,
) (map[uint]abstraction.Asset, error) {
	assets, ok := st.Assets[agentID]
	if !ok {
		return nil, errors.New("Could not find out assets with the agent id")
	}
	return assets, nil
}

func (st *Storage) GetAgentAsset(
	agentID string,
	assetType uint,
) (abstraction.Asset, error) {
	assets, ok := st.Assets[agentID]
	if !ok {
		return nil, errors.New("Could not find out assets with the agent id")
	}
	asset, ok := assets[assetType]
	if !ok {
		return nil, errors.New("Asset type not found.")
	}
	return asset, nil
}

func (st *Storage) GetSortedBidsByAssetType(
	assetType uint,
	isDesc bool,
) abstraction.OrderItems {
	bidsByType, ok := st.Bids[assetType]
	if !ok {
		st.Bids[assetType] = map[string]abstraction.OrderItem{}
		return abstraction.OrderItems{}
	}
	orderItems := abstraction.OrderItems{}
	for _, orderItem := range bidsByType {
		orderItems = append(orderItems, orderItem)
	}
	return orderItems.SortOrderItems(isDesc)
}

func (st *Storage) RemoveBidsByAgentIDs(
	agentIDs []string,
	assetType uint,
) error {
	bidsByAssetType, ok := st.Bids[assetType]
	if !ok {
		return errors.New("Asset type not found.")
	}
	for _, agentID := range agentIDs {
		delete(bidsByAssetType, agentID)
	}
	return nil
}

func (st *Storage) AppendAsk(
	assetType uint,
	agentID string,
	quantity float64,
	pricePerUnit float64,
) {
	orderItem := &OrderItem{
		AgentID:      agentID,
		AssetType:    assetType,
		Quantity:     quantity,
		PricePerUnit: pricePerUnit,
		OrderTime:    time.Now().Unix(),
	}
	asks, ok := st.Asks[assetType]
	if !ok {
		st.Asks[assetType] = map[string]abstraction.OrderItem{
			agentID: orderItem,
		}
		return
	}
	asks[agentID] = orderItem
}

func (st *Storage) GetSortedAsksByAssetType(
	assetType uint,
	isDesc bool,
) abstraction.OrderItems {
	asksByType, ok := st.Asks[assetType]
	if !ok {
		st.Asks[assetType] = map[string]abstraction.OrderItem{}
		return abstraction.OrderItems{}
	}
	orderItems := abstraction.OrderItems{}
	for _, orderItem := range asksByType {
		orderItems = append(orderItems, orderItem)
	}
	return orderItems.SortOrderItems(isDesc)
}

func (st *Storage) RemoveAsksByAgentIDs(
	agentIDs []string,
	assetType uint,
) error {
	asksByAssetType, ok := st.Asks[assetType]
	if !ok {
		return errors.New("Asset type not found.")
	}
	for _, agentID := range agentIDs {
		delete(asksByAssetType, agentID)
	}
	return nil
}

func (st *Storage) AppendBid(
	assetType uint,
	agentID string,
	quantity float64,
	pricePerUnit float64,
) {
	orderItem := &OrderItem{
		AgentID:      agentID,
		AssetType:    assetType,
		Quantity:     quantity,
		PricePerUnit: pricePerUnit,
		OrderTime:    time.Now().Unix(),
	}
	bids, ok := st.Bids[assetType]
	if !ok {
		st.Bids[assetType] = map[string]abstraction.OrderItem{
			agentID: orderItem,
		}
		return
	}
	bids[agentID] = orderItem
}
