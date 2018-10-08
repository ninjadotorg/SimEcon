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
	PERSISTENT_FILE_PATH = "/Users/autonomous/projects/golang-projects/src/github.com/ninjadotorg/SimEcon/agents/miner/persistent.json"
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
	// split balance into 2 parts: one for buying necessity and one for trading or investing
	balForNec := walBal * 0.5
	// balForInvesting := walBal - balForNec

	// get assets
	// agentAssets, err := common.GetAgentAssets(httpClient, agentID)
	// if err != nil {
	// 	fmt.Printf("Get agent assets error: %s\n", err.Error())
	// 	return
	// }

	// buy necessity
	nPricePerUnit := common.NECESSITY_PRICE_BASELINE * ((rand.Float64() * 80) + 40) / 100
	nQty := math.Floor(balForNec / nPricePerUnit)
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

	// TODO: define investment strategy here

	fmt.Printf("Finished the session for agent: %s\n", agentID)
}

func run() {
	httpClient := common.NewHttpClient()
	agentIDs := common.GetAgentIDs(
		httpClient,
		common.GetEnv("PERSISTENT_FILE_PATH", PERSISTENT_FILE_PATH),
		NUMBER_OF_AGENTS,
		common.MINER,
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
