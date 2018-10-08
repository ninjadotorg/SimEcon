package main

import (
	"fmt"
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
	PERSISTENT_FILE_PATH = "/Users/autonomous/projects/golang-projects/src/github.com/ninjadotorg/SimEcon/agents/price-stability/persistent.json"
)

func process(
	httpClient *common.HttpClient,
	agentID string,
) {
	// get coin price
	coinPrice, err := common.GetCoinPrice(httpClient)
	if err != nil {
		fmt.Printf("Get coin price error: %s\n", err.Error())
		return
	}

	// get total tokens
	totalTokens, err := common.GetTotalTokens(httpClient)
	if err != nil {
		fmt.Printf("Get total tokens error: %s\n", err.Error())
		return
	}

	totalCoins := totalTokens.TotalCoins
	totalBonds := totalTokens.TotalBonds
	demand := coinPrice * totalCoins
	delta := demand - totalCoins
	var actionParam *common.ActionParam
	var tax float64 = 0
	if delta > totalBonds {
		tax = 10 // TODO: figure out tax here
	}
	if delta < 0 {
		tax = 0
	}
	actionParam = &common.ActionParam{
		Delta: delta,
		Tax:   tax,
	}
	_, err = common.Stabilize(httpClient, agentID, actionParam)
	if err != nil {
		fmt.Printf("Call stabilization api error: %s\n", err.Error())
		return
	}

	fmt.Printf("Finished the session for agent: %s\n", agentID)
}

func run() {
	httpClient := common.NewHttpClient()
	agentIDs := common.GetAgentIDs(
		httpClient,
		common.GetEnv("PERSISTENT_FILE_PATH", PERSISTENT_FILE_PATH),
		NUMBER_OF_AGENTS,
		common.PRICE_STABILITY,
	)

	// Agent re-calculates every 60s
	deplayTimeInSec, _ := strconv.Atoi(common.GetEnv("DELAY_TIME_IN_SEC", fmt.Sprintf("%d", 5)))
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
