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
	NUMBER_OF_AGENTS     = 20
	PERSISTENT_FILE_PATH = "/Users/autonomous/projects/golang-projects/src/github.com/ninjadotorg/SimEcon/agents/trader/persistent.json"
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
	walCoins := walAcc.Coins
	coinsForTrade := walCoins * 0.6
	coinsForNec := walCoins - coinsForTrade
	walBonds := walAcc.Bonds

	// get assets
	agentAssets, err := common.GetAgentAssets(httpClient, agentID)
	if err != nil {
		fmt.Printf("Get agent assets error: %s\n", err.Error())
		return
	}

	// produce man hours from necessity
	// eat all necessity
	nAsset, _ := agentAssets[common.NECESSITY]
	if nAsset.Quantity > 0 {
		producedAgentAssets, err := common.Produce(
			httpClient,
			agentID,
			map[uint]*common.Asset{common.NECESSITY: nAsset},
		)
		if err != nil {
			fmt.Printf("Produce man hours from necessity error: %s\n", err.Error())
			return
		}

		// sell man hours
		mhAsset, _ := producedAgentAssets[common.MAN_HOUR]
		if mhAsset.Quantity > 0 {
			orderSellItem := &common.OrderItem{
				AgentID:      agentID,
				AssetType:    common.MAN_HOUR,
				Quantity:     mhAsset.Quantity,
				PricePerUnit: common.MAN_HOUR_PRICE_BASELINE * ((rand.Float64() * 80) + 40) / 100,
			}
			_, err = common.Order(
				httpClient,
				agentID,
				orderSellItem,
				"sell",
			)
			if err != nil {
				fmt.Printf("Sell man hours error: %s\n", err.Error())
				return
			}
		}
	} else {
		// get assets
		agentAssets, err = common.GetAgentAssets(httpClient, agentID)
		if err != nil {
			fmt.Printf("Get agent assets error: %s\n", err.Error())
			return
		}
		mhAsset, _ := agentAssets[common.MAN_HOUR]
		if mhAsset.Quantity > 0 {
			orderSellItem := &common.OrderItem{
				AgentID:      agentID,
				AssetType:    common.MAN_HOUR,
				Quantity:     mhAsset.Quantity,
				PricePerUnit: common.MAN_HOUR_PRICE_BASELINE * ((rand.Float64() * 80) + 40) / 100,
			}
			_, err = common.Order(
				httpClient,
				agentID,
				orderSellItem,
				"sell",
			)
			if err != nil {
				fmt.Printf("Sell man hours error: %s\n", err.Error())
				return
			}
		}
	}

	// get coin price
	coinPrice, err := common.GetCoinPrice(httpClient)
	if err != nil {
		fmt.Printf("Get coin price error: %s\n", err.Error())
		return
	}

	if coinPrice > 1.1 && walBonds > 0 { // predict that system will sell coins, trader will buy coins from bonds
		var askingCoinPrice float64 = 1 // TODO: figure out asking coin price here
		qty := walBonds / askingCoinPrice

		orderBuyItem := &common.OrderItem{
			AgentID:      agentID,
			AssetType:    common.COIN,
			Quantity:     qty,
			PricePerUnit: askingCoinPrice,
		}
		_, err = common.Order(
			httpClient,
			agentID,
			orderBuyItem,
			"buyTokens",
		)
		if err != nil {
			fmt.Printf("Buy coins error: %s\n", err.Error())
			return
		}

	} else if coinPrice < 0.9 { // predict that system will sell bonds, trader will buy bonds from coins
		var askingBondPrice float64 = 0.75 // TODO: figure out asking bond price here
		qty := coinsForTrade / askingBondPrice

		orderBuyItem := &common.OrderItem{
			AgentID:      agentID,
			AssetType:    common.BOND,
			Quantity:     qty,
			PricePerUnit: askingBondPrice,
		}
		_, err = common.Order(
			httpClient,
			agentID,
			orderBuyItem,
			"buyTokens",
		)
		if err != nil {
			fmt.Printf("Buy bonds error: %s\n", err.Error())
			return
		}
	}

	// buy necessity
	nPricePerUnit := common.NECESSITY_PRICE_BASELINE * ((rand.Float64() * 80) + 40) / 100
	nQty := math.Floor(coinsForNec / nPricePerUnit)
	if nQty >= 1 {
		orderBuyItem := &common.OrderItem{
			AgentID:      agentID,
			AssetType:    common.NECESSITY,
			Quantity:     nQty,
			PricePerUnit: nPricePerUnit,
		}
		_, err = common.Order(
			httpClient,
			agentID,
			orderBuyItem,
			"buy",
		)
		if err != nil {
			fmt.Printf("Buy necessity error: %s\n", err.Error())
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
		common.PERSON,
	)

	rand_interval := rand.Intn(9) + 2
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
