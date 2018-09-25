package main

import (
	"fmt"
	"math"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/ninjadotorg/SimEcon/common"
)

const (
	NUMBER_OF_AGENTS     = 1
	AGENT_TYPE           = 3
	PERSISTENT_FILE_PATH = "/Users/autonomous/projects/golang-projects/src/github.com/ninjadotorg/SimEcon/micro_economy/cfirm/persistent.json"
)

func process(
	httpClient *common.HttpClient,
	agentID string,
) {
	// get wallet balance
	walBal, err := common.GetWalletBalance(httpClient, agentID)
	if err != nil {
		fmt.Printf("Get wallet balance error: %s\n", err.Error())
		return
	}

	// get assets
	agentAssets, err := common.GetAgentAssets(httpClient, agentID)
	if err != nil {
		fmt.Printf("Get agent assets error: %s\n", err.Error())
		return
	}

	// produce capital from man hours
	mhAsset, _ := agentAssets[common.MAN_HOUR]
	producedAgentAssets, err := common.Produce(
		httpClient,
		agentID,
		map[uint]*common.Asset{common.MAN_HOUR: mhAsset},
	)
	if err != nil {
		fmt.Printf("Produce capital from man hours error: %s\n", err.Error())
		return
	}

	// sell capital
	cAsset, _ := producedAgentAssets[common.CAPITAL]
	orderSellItem := &common.OrderItem{
		AgentID:      agentID,
		AssetType:    common.CAPITAL,
		Quantity:     cAsset.Quantity,
		PricePerUnit: 15,
	}
	_, err = common.Order(
		httpClient,
		agentID,
		orderSellItem,
		"sell",
	)
	if err != nil {
		fmt.Printf("Sell capital error: %s\n", err.Error())
		return
	}

	// buy man hours
	orderBuyItem := &common.OrderItem{
		AgentID:      agentID,
		AssetType:    common.MAN_HOUR,
		Quantity:     math.Floor(walBal / 21),
		PricePerUnit: 21,
	}
	_, err = common.Order(
		httpClient,
		agentID,
		orderBuyItem,
		"buy",
	)
	if err != nil {
		fmt.Printf("Buy man hours error: %s\n", err.Error())
		return
	}

	fmt.Println("Everything is ok")
}

func run() {
	httpClient := common.NewHttpClient()
	agentIDs := common.GetAgentIDs(
		httpClient,
		common.GetEnv("PERSISTENT_FILE_PATH", PERSISTENT_FILE_PATH),
		NUMBER_OF_AGENTS,
		AGENT_TYPE,
	)

	// Agent re-calculates every 60s
	deplayTimeInSec, _ := strconv.Atoi(common.GetEnv("DELAY_TIME_IN_SEC", "600"))
	for {
		fmt.Println("Hello there again!!!")
		for _, agentID := range agentIDs {
			go process(httpClient, agentID)
		}

		time.Sleep(time.Duration(deplayTimeInSec) * time.Second)
	}
}

func clearUpBeforeTerminating() {
	// TODO: do cleaning up here, probably send message to a channel in run func to stop loops
	fmt.Println("Wait for 2 seconds to finish processing")
	time.Sleep(2 * time.Second)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	var gracefulStop = make(chan os.Signal)
	signal.Notify(gracefulStop, syscall.SIGTERM)
	signal.Notify(gracefulStop, syscall.SIGINT)
	go func() {
		sig := <-gracefulStop
		fmt.Printf("caught sig: %+v", sig)
		clearUpBeforeTerminating()
		os.Exit(0)
	}()
	run()
}
