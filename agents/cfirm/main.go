package main

import (
	"fmt"
	"math"
	"math/rand"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/ninjadotorg/SimEcon/common"
)

const (
	NUMBER_OF_AGENTS     = 5
	AGENT_TYPE           = 3
	PERSISTENT_FILE_PATH = "/Users/autonomous/projects/golang-projects/src/github.com/ninjadotorg/SimEcon/agents/cfirm/persistent.json"
)

func process(
	httpClient *common.HttpClient,
	agentID string,
) {
	// get wallet account
	walAcc, err := common.GetWalletAccount(httpClient, agentID)
	if err != nil {
		fmt.Printf("Get wallet balance error: %s\n", err.Error())
		return
	}
	walBal := walAcc.Coins

	// get assets
	agentAssets, err := common.GetAgentAssets(httpClient, agentID)
	if err != nil {
		fmt.Printf("Get agent assets error: %s\n", err.Error())
		return
	}

	// produce capital from man hours
	mhAsset, _ := agentAssets[common.MAN_HOUR]
	if mhAsset.Quantity > 0 {
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
			PricePerUnit: common.CAPITAL_PRICE_BASELINE * ((rand.Float64() * 80) + 40) / 100,
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
	} else {
		// get assets
		agentAssets, err := common.GetAgentAssets(httpClient, agentID)
		if err != nil {
			fmt.Printf("Get agent assets error: %s\n", err.Error())
			return
		}
		// sell capital
		cAsset, _ := agentAssets[common.CAPITAL]
		if cAsset.Quantity > 0 {
			orderSellItem := &common.OrderItem{
				AgentID:      agentID,
				AssetType:    common.CAPITAL,
				Quantity:     cAsset.Quantity,
				PricePerUnit: common.CAPITAL_PRICE_BASELINE * ((rand.Float64() * 80) + 40) / 100,
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
		}
	}

	// buy man hours
	mhPricePerUnit := common.MAN_HOUR_PRICE_BASELINE * ((rand.Float64() * 80) + 40) / 100
	mhQty := math.Floor(walBal / mhPricePerUnit)
	if mhQty >= 1 {
		orderBuyItem := &common.OrderItem{
			AgentID:      agentID,
			AssetType:    common.MAN_HOUR,
			Quantity:     mhQty,
			PricePerUnit: mhPricePerUnit,
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
	}

	fmt.Printf("Finished the session for agent: %s\n", agentID)
}

func run() {
	httpClient := common.NewHttpClient()
	agentIDs := common.GetAgentIDs(
		httpClient,
		common.GetEnv("PERSISTENT_FILE_PATH", PERSISTENT_FILE_PATH),
		NUMBER_OF_AGENTS,
		AGENT_TYPE,
	)

	rand_interval := rand.Intn(8) + 2
	// rand_interval := 600
	// Agent re-calculates every 60s
	deplayTimeInSec, _ := strconv.Atoi(common.GetEnv("DELAY_TIME_IN_SEC", fmt.Sprintf("%d", rand_interval)))
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
