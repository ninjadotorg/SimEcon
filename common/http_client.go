package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

type HttpClient struct {
	*http.Client
}

func NewHttpClient() *HttpClient {
	httpClient := &http.Client{
		Timeout: time.Second * 60,
	}
	return &HttpClient{
		httpClient,
	}
}

func BuildHttpServerAddress() string {
	protocol := GetEnv("PROTOCOL", "http")
	host := GetEnv("HOST", "127.0.0.1")
	port, _ := strconv.Atoi(GetEnv("PORT", "8080"))
	return fmt.Sprintf("%s://%s:%d", protocol, host, port)
}

func BuildJoinSimulationEndPoint(agentType uint) string {
	return fmt.Sprintf("%s/types/%d/agents", BuildHttpServerAddress(), agentType)
}

func BuildGetWalletBalanceEndPoint(agentID string) string {
	return fmt.Sprintf("%s/agents/%s/wallet/balance", BuildHttpServerAddress(), agentID)
}

func BuildGetAgentAssetsEndPoint(agentID string) string {
	return fmt.Sprintf("%s/agents/%s/assets", BuildHttpServerAddress(), agentID)
}

func BuildProduceEndPoint(agentID string) string {
	return fmt.Sprintf("%s/agents/%s/produce", BuildHttpServerAddress(), agentID)
}

func BuildBuyEndPoint(agentID string) string {
	return fmt.Sprintf("%s/agents/%s/buy", BuildHttpServerAddress(), agentID)
}

func BuildSellEndPoint(agentID string) string {
	return fmt.Sprintf("%s/agents/%s/sell", BuildHttpServerAddress(), agentID)
}

func HandleHttpResp(
	result interface{},
	httpResp *http.Response,
	err error,
) error {
	if err != nil {
		return err
	}
	respBody := httpResp.Body
	defer respBody.Close()

	body, err := ioutil.ReadAll(respBody)
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, result)
	if err != nil {
		return err
	}
	return nil
}
