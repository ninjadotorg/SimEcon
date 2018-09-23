package main

import (
	"github.com/ninjadotorg/SimEcon/macro_economy/account_manager"
	"github.com/ninjadotorg/SimEcon/macro_economy/economy"
	"github.com/ninjadotorg/SimEcon/macro_economy/market"
	"github.com/ninjadotorg/SimEcon/macro_economy/production"
	"github.com/ninjadotorg/SimEcon/macro_economy/storage"
)

func main() {
	st := storage.GetStorageInstance()
	ac := account_manager.GetAccountManagerInstance()
	prod := production.GetProductionInstance()
	m := market.GetMarketInstance()
	econ := economy.GetEconomyInstance(ac, st, prod, m)
	econ.Run()
}
