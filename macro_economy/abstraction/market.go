package abstraction

import "github.com/ninjadotorg/SimEcon/macro_economy/dto"

type Market interface {
	Buy(string, *dto.OrderItem, Storage, AccountManager, Production, Tracker) (float64, error)
	Sell(string, *dto.OrderItem, Storage, AccountManager, Production, Tracker) (float64, error)
	BuyTokens(string, *dto.OrderItem, Storage, AccountManager, Tracker) (float64, error)
	SellTokens(*dto.OrderItem, Storage, AccountManager, Tracker) (float64, error)
}
