package storage

import (
	"fmt"
	"time"

	"github.com/ninjadotorg/SimEcon/common"
	"github.com/ninjadotorg/SimEcon/macro_economy/abstraction"
	"github.com/ninjadotorg/SimEcon/macro_economy/dto"
)

type Person struct {

	// // savings rate (portion of total income+savings that is saved in the last step)
	// SavingsRate float64

	// // consumption (in coin)
	// Consumption float64

	// // minimum necessity (in real quantity) to buy in the current step
	// MinN float64

	// // lowest real interest rate seen
	// LowRR float64

	// // highest real interest rate seen
	// HighRR float64

	Agent
}

func (p *Person) InitAgentAssets(
	st abstraction.Storage,
) {
	// necessity asset
	nAsset := &Asset{
		AgentID:      p.AgentID,
		Type:         common.NECESSITY,
		Quantity:     common.PERSON_NECESSITY,
		ProducedTime: time.Now().Unix(),
	}

	mhAsset := &Asset{
		AgentID:      p.AgentID,
		Type:         common.MAN_HOUR,
		Quantity:     common.PERSON_MAN_HOURS,
		ProducedTime: time.Now().Unix(),
	}
	st.UpdateAssets(p.AgentID, []abstraction.Asset{nAsset, mhAsset})
}

func (p *Person) GetType() uint {
	return p.Type
}

func (p *Person) UpdateAgent(
	st abstraction.Storage,
	agentDTO abstraction.AgentDTO,
) {
	updatingPerson := agentDTO.(dto.Person)
	fmt.Println("updatingPerson: ", updatingPerson)
}
