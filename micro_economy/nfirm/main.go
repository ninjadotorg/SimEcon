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
	NUMBER_OF_AGENTS     = 3
	AGENT_TYPE           = 2
	PERSISTENT_FILE_PATH = "/Users/autonomous/projects/golang-projects/src/github.com/ninjadotorg/SimEcon/micro_economy/nfirm/persistent.json"
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

	// produce necessity from capital & man hours
	mhAsset, _ := agentAssets[common.MAN_HOUR]
	cAsset, _ := agentAssets[common.CAPITAL]
	producedAgentAssets, err := common.Produce(
		httpClient,
		agentID,
		map[uint]*common.Asset{
			common.MAN_HOUR: mhAsset,
			common.CAPITAL:  cAsset,
		},
	)
	if err != nil {
		fmt.Printf("Produce necessity from capital & man hours error: %s\n", err.Error())
		return
	}

	// sell necessity
	nAsset, _ := producedAgentAssets[common.NECESSITY]
	orderSellItem := &common.OrderItem{
		AgentID:      agentID,
		AssetType:    common.NECESSITY,
		Quantity:     nAsset.Quantity,
		PricePerUnit: 12,
	}
	_, err = common.Order(
		httpClient,
		agentID,
		orderSellItem,
		"sell",
	)
	if err != nil {
		fmt.Printf("Sell necessity error: %s\n", err.Error())
		return
	}

	// buy capital
	pricePerCapital := 9.0
	balForCapital := walBal / 2
	orderBuyCapital := &common.OrderItem{
		AgentID:      agentID,
		AssetType:    common.CAPITAL,
		Quantity:     math.Floor(balForCapital / pricePerCapital),
		PricePerUnit: pricePerCapital,
	}
	_, err = common.Order(
		httpClient,
		agentID,
		orderBuyCapital,
		"buy",
	)
	if err != nil {
		fmt.Printf("Buy capital error: %s\n", err.Error())
		return
	}

	// buy man hours
	pricePerManHour := 22.0
	balForManHour := walBal - balForCapital
	orderBuyManHour := &common.OrderItem{
		AgentID:      agentID,
		AssetType:    common.CAPITAL,
		Quantity:     math.Floor(balForManHour / pricePerManHour),
		PricePerUnit: pricePerManHour,
	}
	_, err = common.Order(
		httpClient,
		agentID,
		orderBuyManHour,
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