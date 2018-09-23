package economy

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ninjadotorg/SimEcon/common"
	"github.com/ninjadotorg/SimEcon/macro_economy/dto"
)

// POST /types/{AGENT_TYPE}/agents
func Join(w http.ResponseWriter, r *http.Request, econ *Economy) {
	newAgentID := common.UUID()
	am := econ.AccountManager
	st := econ.Storage
	agentType, _ := strconv.Atoi(mux.Vars(r)["AGENT_TYPE"])

	// open wallet account
	am.OpenWalletAccount(newAgentID, 0.0)

	// insert new agent
	agent := st.InsertAgent(newAgentID, uint(agentType))
	agent.InitAgentAssets(st)

	jsInBytes, _ := json.Marshal(agent)
	w.Write(jsInBytes)
}

// POST /agents/{AGENT_ID}/produce
func Produce(w http.ResponseWriter, r *http.Request, econ *Economy) {
	st := econ.Storage
	prod := econ.Production

	var assetsReq map[uint]*dto.Asset
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&assetsReq)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	// TODO: validate assets requested by agent type

	agentID := mux.Vars(r)["AGENT_ID"]
	agent, _ := st.GetAgentByID(agentID)

	agentProd, err := prod.GetProductionByAgentType(agent.GetType())
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}
	updatedAssets, err := agentProd.Produce(st, agentID, assetsReq)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	jsInBytes, _ := json.Marshal(updatedAssets)
	w.Write(jsInBytes)
}

// GET /agents/{AGENT_ID}/assets
func GetAgentAssets(w http.ResponseWriter, r *http.Request, econ *Economy) {
	agentID := mux.Vars(r)["AGENT_ID"]
	st := econ.Storage
	assets, err := st.GetAgentAssets(agentID)

	// TODO: get actual qty of these assets

	res := map[string]interface{}{}
	if err != nil {
		res["error"] = err.Error()
	} else {
		res["result"] = assets
	}
	jsInBytes, _ := json.Marshal(res)
	w.Write(jsInBytes)
}

// POST /agents/{AGENT_ID}/buy
func Buy(w http.ResponseWriter, r *http.Request, econ *Economy) {
	st := econ.Storage
	mk := econ.Market
	am := econ.AccountManager
	prod := econ.Production

	agentID := mux.Vars(r)["AGENT_ID"]
	var orderItemReq *dto.OrderItem
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&orderItemReq)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	agent, err := st.GetAgentByID(agentID)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	agentAsset, err := st.GetAgentAsset(agentID, orderItemReq.AssetType)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	agentProd, err := prod.GetProductionByAgentType(agent.GetType())
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	curAsset := agentProd.GetActualAsset(agentAsset)
	accBal := am.GetBalance(agentID)
	res := map[string]interface{}{
		"message":           "Process buy order successfully.",
		"oldAssetQuantity":  curAsset.GetQuantity(),
		"oldAccountBalance": accBal,
	}

	reqQty := orderItemReq.Quantity

	// validate if coin balance is enough for the order or not?
	// TODO: should validate the other buy orders of this agent (on other asset type)
	if accBal < orderItemReq.Quantity*orderItemReq.PricePerUnit {
		res = map[string]interface{}{
			"message":        "Not enough money for the buy order",
			"accountBalance": accBal,
			"orderQuantity":  reqQty,
			"pricePerUnit":   orderItemReq.PricePerUnit,
		}
		jsInBytes, _ := json.Marshal(res)
		w.Write(jsInBytes)
		return
	}

	remainingRequestedQty, err := mk.Buy(agentID, orderItemReq, st, am)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	newAssetQuantity := curAsset.GetQuantity() + reqQty - remainingRequestedQty
	// update asset after buying
	curAsset.SetQuantity(newAssetQuantity)
	st.UpdateAsset(agentID, curAsset)

	res["newAccountBalance"] = am.GetBalance(agentID)
	res["newAssetQuantity"] = newAssetQuantity
	jsInBytes, _ := json.Marshal(res)
	w.Write(jsInBytes)
}

// POST /agents/{AGENT_ID}/sell
func Sell(w http.ResponseWriter, r *http.Request, econ *Economy) {
	st := econ.Storage
	mk := econ.Market
	am := econ.AccountManager
	prod := econ.Production

	agentID := mux.Vars(r)["AGENT_ID"]
	var orderItemReq *dto.OrderItem
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&orderItemReq)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	agent, err := st.GetAgentByID(agentID)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	agentAsset, err := st.GetAgentAsset(agentID, orderItemReq.AssetType)
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	agentProd, err := prod.GetProductionByAgentType(agent.GetType())
	if err != nil {
		w.Write([]byte(err.Error()))
		return
	}

	reqQty := orderItemReq.Quantity
	curAsset := agentProd.GetActualAsset(agentAsset)
	if curAsset.GetQuantity() < orderItemReq.Quantity {
		res := map[string]interface{}{
			"message":           "Not enough asset quantity for the sell order",
			"assetType":         orderItemReq.AssetType,
			"requestedQuantity": reqQty,
			"remainingQuantity": curAsset.GetQuantity(),
		}
		jsInBytes, _ := json.Marshal(res)
		w.Write(jsInBytes)
		return
	}

	accBal := am.GetBalance(agentID)
	res := map[string]interface{}{
		"message":           "Process sell order successfully.",
		"oldAssetQuantity":  curAsset.GetQuantity(),
		"oldAccountBalance": accBal,
	}
	remainingRequestedQty, err := mk.Sell(agentID, orderItemReq, st, am)
	if err != nil {
		w.Write([]byte(err.Error()))
	}

	newAssetQuantity := curAsset.GetQuantity() - (reqQty - remainingRequestedQty)
	// update asset after buying
	curAsset.SetQuantity(newAssetQuantity)
	st.UpdateAsset(agentID, curAsset)

	res["newAccountBalance"] = am.GetBalance(agentID)
	res["newAssetQuantity"] = newAssetQuantity
	jsInBytes, _ := json.Marshal(res)
	w.Write(jsInBytes)
}
