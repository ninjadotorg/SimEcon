package abstraction

import (
	"github.com/ninjadotorg/SimEcon/macro_economy/dto"
)

type AgentProduction interface {
	Produce(Storage, string, map[uint]*dto.Asset) (map[uint]Asset, error)
}
