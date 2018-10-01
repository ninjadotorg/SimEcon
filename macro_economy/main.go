package main

import (
	"github.com/ninjadotorg/SimEcon/macro_economy/account_manager"
	"github.com/ninjadotorg/SimEcon/macro_economy/economy"
	"github.com/ninjadotorg/SimEcon/macro_economy/market"
	"github.com/ninjadotorg/SimEcon/macro_economy/production"
	"github.com/ninjadotorg/SimEcon/macro_economy/storage"
	"github.com/ninjadotorg/SimEcon/macro_economy/tracker"
)

func main() {
	st := storage.GetStorageInstance()
	ac := account_manager.GetAccountManagerInstance()
	prod := production.GetProductionInstance()
	m := market.GetMarketInstance()
	tr := tracker.GetTrackerInstance()
	econ := economy.GetEconomyInstance(ac, st, prod, m, tr)
	econ.Run()
}
