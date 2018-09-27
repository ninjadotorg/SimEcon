package economy

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/ninjadotorg/SimEcon/common"
	"github.com/ninjadotorg/SimEcon/macro_economy/abstraction"
	"github.com/ninjadotorg/SimEcon/macro_economy/dto"
)

// TODO: will remove unnecessary mutex locks when using real DB (Redis, Couchbase, LevelDB, ...)

var counter = 0

func validateRequestedProduce(agentType uint, assetsReq map[uint]*dto.Asset) bool {
	assetTypesReq := []uint{}
	for assetType, _ := range assetsReq {
		assetTypesReq = append(assetTypesReq, assetType)
	}
	validationMap := map[uint][]uint{
		common.PERSON:         []uint{common.NECESSITY},
		common.NECESSITY_FIRM: []uint{common.MAN_HOUR, common.CAPITAL},
		common.CAPITAL_FIRM:   []uint{common.MAN_HOUR},
	}
	if len(validationMap[agentType]) != len(assetTypesReq) {
		return false
	}
	for _, validatedAssetType := range validationMap[agentType] {
		found := false
		for _, assetTypeReq := range assetTypesReq {
			if assetTypeReq == validatedAssetType {
				found = true
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func decodePersonDTO(
	r *http.Request,
) (abstraction.AgentDTO, error) {
	var personDTO dto.Person
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&personDTO)
	if err != nil {
		return nil, err
	}
	return personDTO, nil
}

func decodeNecessityFirmDTO(
	r *http.Request,
) (abstraction.AgentDTO, error) {
	var nFirmDTO dto.NecessityFirm
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&nFirmDTO)
	if err != nil {
		return nil, err
	}
	return nFirmDTO, nil
}

func decodeCapitalFirmDTO(
	r *http.Request,
) (abstraction.AgentDTO, error) {
	var cFirmDTO dto.CapitalFirm
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&cFirmDTO)
	if err != nil {
		return nil, err
	}
	return cFirmDTO, nil
}

// POST /types/{AGENT_TYPE}/agents
func Join(w http.ResponseWriter, r *http.Request, econ *Economy) {
	counter += 1
	var mutex = &sync.Mutex{}
	newAgentID := common.UUID()

	am := econ.AccountManager
	st := econ.Storage
	agentType, err := strconv.Atoi(mux.Vars(r)["AGENT_TYPE"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// //just for demo
	// if counter == 1 {
	// 	newAgentID = "04103203-829E-452D-B9FE-A8017F7541B6"
	// }
	// if counter == 2 {
	// 	newAgentID = "AE8B49DF-B96E-4E87-AAB9-B7B7E16B48D0"
	// }
	// if counter == 3 {
	// 	newAgentID = "E44DAF38-D2E8-4DBC-9A37-29F917A7DB0F"
	// }

	mutex.Lock()
	// open wallet account
	var initBal float64 = common.DEFAULT_ACCOUNT_BALANCE
	if agentType == common.NECESSITY_FIRM {
		initBal = common.NFIRM_ACCOUNT_BALANCE
	}
	am.OpenWalletAccount(newAgentID, initBal)

	// insert new agent
	agent := st.InsertAgent(newAgentID, uint(agentType))
	agent.InitAgentAssets(st)
	mutex.Unlock()

	jsInBytes, _ := json.Marshal(agent)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsInBytes)
}

// GET /agents/{AGENT_ID}
func GetAgentByID(w http.ResponseWriter, r *http.Request, econ *Economy) {
	st := econ.Storage
	agentID := mux.Vars(r)["AGENT_ID"]
	agent, err := st.GetAgentByID(agentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	jsInBytes, _ := json.Marshal(agent)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsInBytes)
}

// PUT /agents/{AGENT_ID}
func UpdateAgent(w http.ResponseWriter, r *http.Request, econ *Economy) {
	st := econ.Storage
	agentID := mux.Vars(r)["AGENT_ID"]
	agent, err := st.GetAgentByID(agentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	agentType := agent.GetType()
	var agentDTO abstraction.AgentDTO
	if agentType == common.PERSON {
		agentDTO, err = decodePersonDTO(r)
	} else if agentType == common.NECESSITY_FIRM {
		agentDTO, err = decodeNecessityFirmDTO(r)
	} else if agentType == common.CAPITAL_FIRM {
		agentDTO, err = decodeCapitalFirmDTO(r)
	}

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	agent.UpdateAgent(st, agentDTO)
}

// POST /agents/{AGENT_ID}/produce
func Produce(w http.ResponseWriter, r *http.Request, econ *Economy) {
	var mutex = &sync.Mutex{}
	st := econ.Storage
	prod := econ.Production

	var assetsReq map[uint]*dto.Asset
	decoder := json.NewDecoder(r.Body)

	err := decoder.Decode(&assetsReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	agentID := mux.Vars(r)["AGENT_ID"]
	agent, err := st.GetAgentByID(agentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	isValid := validateRequestedProduce(agent.GetType(), assetsReq)
	if !isValid {
		http.Error(w, "Assets requested are invalid for this agent type.", http.StatusBadRequest)
		return
	}

	agentProd, err := prod.GetProductionByAgentType(agent.GetType())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	mutex.Lock()
	updatedAssets, err := agentProd.Produce(st, agentID, assetsReq)
	mutex.Unlock()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	jsInBytes, _ := json.Marshal(updatedAssets)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsInBytes)
}

// GET /agents/{AGENT_ID}/assets
func GetAgentAssets(w http.ResponseWriter, r *http.Request, econ *Economy) {
	agentID := mux.Vars(r)["AGENT_ID"]
	st := econ.Storage
	prod := econ.Production
	assets, err := st.GetAgentAssets(agentID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	actualAssets := prod.GetActualAssets(assets)
	jsInBytes, _ := json.Marshal(actualAssets)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsInBytes)
}

// GET /agents/{AGENT_ID}/wallet/balance
func GetWalletAccountBalance(w http.ResponseWriter, r *http.Request, econ *Economy) {
	am := econ.AccountManager
	agentID := mux.Vars(r)["AGENT_ID"]
	bal := am.GetBalance(agentID)
	jsInBytes, _ := json.Marshal(
		map[string]interface{}{
			"agentId": agentID,
			"balance": bal,
		},
	)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsInBytes)
}

// GET /agents/{AGENT_ID}/wallet/account
func GetWalletAccount(w http.ResponseWriter, r *http.Request, econ *Economy) {
	am := econ.AccountManager
	agentID := mux.Vars(r)["AGENT_ID"]
	walletAcc := am.GetWalletAccount(agentID)
	jsInBytes, _ := json.Marshal(
		map[string]interface{}{
			"agentId":       agentID,
			"walletAccount": walletAcc,
		},
	)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsInBytes)
}

// POST /agents/{AGENT_ID}/buy
func Buy(w http.ResponseWriter, r *http.Request, econ *Economy) {
	var mutex = &sync.Mutex{}
	st := econ.Storage
	mk := econ.Market
	am := econ.AccountManager
	prod := econ.Production

	agentID := mux.Vars(r)["AGENT_ID"]
	var orderItemReq *dto.OrderItem
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&orderItemReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	agentAsset, err := st.GetAgentAsset(agentID, orderItemReq.AssetType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	curAsset := prod.GetActualAsset(agentAsset)
	accBal := am.GetBalance(agentID)
	reqQty := orderItemReq.Quantity

	// validate if coin balance is enough for the order or not?
	// TODO: should validate the other buy orders of this agent (on other asset type)
	if accBal < orderItemReq.Quantity*orderItemReq.PricePerUnit {
		res := map[string]interface{}{
			"message":        "Not enough money for the buy order",
			"accountBalance": accBal,
			"orderQuantity":  reqQty,
			"pricePerUnit":   orderItemReq.PricePerUnit,
		}
		jsInBytes, _ := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsInBytes)
		return
	}

	res := map[string]interface{}{
		"message":           "Process buy order successfully.",
		"oldAssetQuantity":  curAsset.GetQuantity(),
		"oldAccountBalance": accBal,
	}

	mutex.Lock()
	remainingRequestedQty, err := mk.Buy(agentID, orderItemReq, st, am, prod)
	mutex.Unlock()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newAssetQuantity := curAsset.GetQuantity() + reqQty - remainingRequestedQty

	// update asset after buying
	curAsset.SetQuantity(newAssetQuantity)
	st.UpdateAsset(agentID, curAsset)

	res["newAccountBalance"] = am.GetBalance(agentID)
	res["newAssetQuantity"] = newAssetQuantity
	jsInBytes, _ := json.Marshal(res)
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsInBytes)
}

// POST /agents/{AGENT_ID}/sell
func Sell(w http.ResponseWriter, r *http.Request, econ *Economy) {
	var mutex = &sync.Mutex{}
	st := econ.Storage
	mk := econ.Market
	am := econ.AccountManager
	prod := econ.Production

	agentID := mux.Vars(r)["AGENT_ID"]
	var orderItemReq *dto.OrderItem
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&orderItemReq)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	agentAsset, err := st.GetAgentAsset(agentID, orderItemReq.AssetType)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	reqQty := orderItemReq.Quantity
	curAsset := prod.GetActualAsset(agentAsset)
	if curAsset.GetQuantity() < orderItemReq.Quantity {
		res := map[string]interface{}{
			"message":           "Not enough asset quantity for the sell order",
			"assetType":         orderItemReq.AssetType,
			"requestedQuantity": reqQty,
			"remainingQuantity": curAsset.GetQuantity(),
		}
		jsInBytes, _ := json.Marshal(res)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsInBytes)
		return
	}

	accBal := am.GetBalance(agentID)
	res := map[string]interface{}{
		"message":           "Process sell order successfully.",
		"oldAssetQuantity":  curAsset.GetQuantity(),
		"oldAccountBalance": accBal,
	}
	mutex.Lock()
	remainingRequestedQty, err := mk.Sell(agentID, orderItemReq, st, am, prod)
	mutex.Unlock()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	newAssetQuantity := curAsset.GetQuantity() - (reqQty - remainingRequestedQty)
	// update asset after buying
	curAsset.SetQuantity(newAssetQuantity)
	st.UpdateAsset(agentID, curAsset)

	res["newAccountBalance"] = am.GetBalance(agentID)
	res["newAssetQuantity"] = newAssetQuantity
	jsInBytes, _ := json.Marshal(res)

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsInBytes)
}
