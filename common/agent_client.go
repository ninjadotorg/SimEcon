package common

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

type ActionParam struct {
	Delta float64 `json:"delta"`
	Tax   float64 `json:"tax"`
}

type PersistedContent struct {
	AgentIDs []string `json:"agentIds"`
}

type Asset struct {
	AgentID      string  `json:"AgentID"`
	Type         uint    `json:"Type"`
	Quantity     float64 `json:"Quantity"`
	ProducedTime int64   `json:"ProducedTime"`
}

type OrderItem struct {
	AgentID      string  `json:"AgentID"`
	AssetType    uint    `json:"AssetType"`
	Quantity     float64 `json:"Quantity"`
	PricePerUnit float64 `json:"PricePerUnit"`
}

type OrderResponse struct {
	Message           string  `json:"message"`
	NewAccountBalance float64 `json:"newAccountBalance"`
	NewAssetQuantity  float64 `json:"newAssetQuantity"`
	OldAccountBalance float64 `json:"oldAccountBalance"`
	OldAssetQuantity  float64 `json:"oldAssetQuantity"`
}

type WalletAccount struct {
	Address string
	Coins   float64
	Bonds   float64
	PriIC   float64 // primary income in the last step
	SecIC   float64 // secondary income in the last step
}

type WalletAccResp struct {
	AgentID       string         `json:"agentId"`
	WalletAccount *WalletAccount `json:"walletAccount"`
}

type CoinPriceResp struct {
	CoinPrice float64 `json:"coinPrice"`
}

type TotalTokensResp struct {
	TotalCoins float64 `json:"totalCoins"`
	TotalBonds float64 `json:"totalBonds"`
}

func GetTotalTokens(
	httpClient *HttpClient,
) (*TotalTokensResp, error) {
	resp, err := httpClient.Get(BuildGetTotalTokensEndPoint())
	if err != nil {
		return nil, err
	}
	var totalTokensResp TotalTokensResp
	err = HandleHttpResp(&totalTokensResp, resp, err)
	if err != nil {
		return nil, err
	}
	return &totalTokensResp, nil
}

func GetCoinPrice(
	httpClient *HttpClient,
) (float64, error) {
	resp, err := httpClient.Get(BuildGetCoinPriceEndPoint())
	if err != nil {
		return 0, err
	}
	var coinPriceResp CoinPriceResp
	err = HandleHttpResp(&coinPriceResp, resp, err)
	if err != nil {
		return 0, err
	}
	return coinPriceResp.CoinPrice, nil
}

func GetWalletAccount(
	httpClient *HttpClient,
	agentID string,
) (*WalletAccount, error) {
	resp, err := httpClient.Get(BuildGetWalletAccountEndPoint(agentID))
	if err != nil {
		return nil, err
	}
	var walletAccResp WalletAccResp
	err = HandleHttpResp(&walletAccResp, resp, err)
	if err != nil {
		return nil, err
	}
	return walletAccResp.WalletAccount, nil
}

func GetWalletBalance(
	httpClient *HttpClient,
	agentID string,
) (float64, error) {
	resp, err := httpClient.Get(BuildGetWalletBalanceEndPoint(agentID))
	if err != nil {
		return 0, err
	}
	var walletBalResp struct {
		AgentID string  `json:"agentId"`
		Balance float64 `json:"balance"`
	}
	err = HandleHttpResp(&walletBalResp, resp, err)
	if err != nil {
		return 0, err
	}
	return walletBalResp.Balance, nil
}

func GetAgentAssets(
	httpClient *HttpClient,
	agentID string,
) (map[uint]*Asset, error) {
	resp, err := httpClient.Get(BuildGetAgentAssetsEndPoint(agentID))
	if err != nil {
		return nil, err
	}
	var agentAssetsResp map[uint]*Asset
	err = HandleHttpResp(&agentAssetsResp, resp, err)
	if err != nil {
		return nil, err
	}
	return agentAssetsResp, nil
}

func Produce(
	httpClient *HttpClient,
	agentID string,
	agentAssets map[uint]*Asset,
) (map[uint]*Asset, error) {
	payloadInBytes, err := json.Marshal(agentAssets)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Post(
		BuildProduceEndPoint(agentID),
		"application/json",
		bytes.NewBuffer(payloadInBytes),
	)
	if err != nil {
		return nil, err
	}
	var agentAssetsResp map[uint]*Asset
	err = HandleHttpResp(&agentAssetsResp, resp, err)
	if err != nil {
		return nil, err
	}
	return agentAssetsResp, nil
}

func Order(
	httpClient *HttpClient,
	agentID string,
	orderItem *OrderItem,
	orderType string,
) (*OrderResponse, error) {
	orderMap := map[string]func(string) string{
		"buy":       BuildBuyEndPoint,
		"sell":      BuildSellEndPoint,
		"buyTokens": BuildBuyTokensEndPoint,
	}
	payloadInBytes, err := json.Marshal(orderItem)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Post(
		orderMap[orderType](agentID),
		"application/json",
		bytes.NewBuffer(payloadInBytes),
	)
	if err != nil {
		return nil, err
	}
	var orderResponse *OrderResponse
	err = HandleHttpResp(&orderResponse, resp, err)
	if err != nil {
		return nil, err
	}
	return orderResponse, nil
}

func GetAgentIDs(
	httpClient *HttpClient,
	persistentFile string,
	numberOfAgents int,
	agentType uint,
) []string {
	var agentIDs = []string{}
	// load agentIDs list from file
	jsonFile, err := os.Open(persistentFile)
	if err != nil {
		fmt.Printf("Open file error: %s\n", err.Error())
	} else {
		byteValue, _ := ioutil.ReadAll(jsonFile)
		jsonFile.Close()

		var persistedContent PersistedContent
		json.Unmarshal(byteValue, &persistedContent)
		agentIDs = persistedContent.AgentIDs
	}

	if len(agentIDs) == 0 {
		// TODO: might split into chunks to spped up this process
		for i := 0; i < numberOfAgents; i++ {
			resp, err := httpClient.Post(
				BuildJoinSimulationEndPoint(agentType),
				"application/json",
				bytes.NewBuffer([]byte{}),
			)
			var result struct {
				AgentID string `json:"AgentID"`
				Type    uint   `json:"Type"`
			}
			err = HandleHttpResp(&result, resp, err)
			if err != nil {
				fmt.Printf("Join simulation error: %s\n", err.Error())
				continue
			}
			agentIDs = append(agentIDs, result.AgentID)
		}
		// persist agent ids to file
		fileContent, _ := json.Marshal(
			map[string]interface{}{
				"agentIds": agentIDs,
			},
		)
		ioutil.WriteFile(persistentFile, fileContent, 0644)
	}
	return agentIDs
}

func Stabilize(
	httpClient *HttpClient,
	agentID string,
	actionParam *ActionParam,
) (*OrderResponse, error) {
	payloadInBytes, err := json.Marshal(actionParam)
	if err != nil {
		return nil, err
	}
	resp, err := httpClient.Post(
		BuildStabilizeEndPoint(agentID),
		"application/json",
		bytes.NewBuffer(payloadInBytes),
	)
	if err != nil {
		return nil, err
	}

	var orderResponse *OrderResponse
	err = HandleHttpResp(&orderResponse, resp, err)
	if err != nil {
		return nil, err
	}
	return orderResponse, nil
}
