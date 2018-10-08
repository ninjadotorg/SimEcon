package abstraction

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/ninjadotorg/SimEcon/common"
)

type Abstraction interface {
}

type Agent struct {
	AgentType          uint
	Num                int
	PersistentFilePath string
}

func (a *Agent) Run(
	isRandRun bool,
	randRunMin int,
	distance int,
	process func(*common.HttpClient, string),
) {
	httpClient := common.NewHttpClient()
	agentIDs := common.GetAgentIDs(
		httpClient,
		a.PersistentFilePath,
		a.Num,
		a.AgentType,
	)
	intervalRun := 0
	if isRandRun {
		intervalRun = rand.Intn(distance) + randRunMin
	}
	deplayTimeInSec, _ := strconv.Atoi(common.GetEnv("DELAY_TIME_IN_SEC", fmt.Sprintf("%d", intervalRun)))
	for {
		fmt.Println("Hello there again!!!")
		for _, agentID := range agentIDs {
			go process(httpClient, agentID)
		}
		time.Sleep(time.Duration(deplayTimeInSec) * time.Second)
	}
}

func name() {

}
