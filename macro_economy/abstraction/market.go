package abstraction

import "github.com/ninjadotorg/SimEcon/macro_economy/dto"

type Market interface {
	Buy(string, *dto.OrderItem, Storage, AccountManager, Production, Tracker) (float64, error)
	Sell(string, *dto.OrderItem, Storage, AccountManager, Production, Tracker) (float64, error)
}
