package economy

import (
	"fmt"
	"math/rand"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/ninjadotorg/SimEcon/common"
	"github.com/ninjadotorg/SimEcon/macro_economy/abstraction"
)

const (
	// SERVER
	PORT = ":8080"
)

type Economy struct {
	AccountManager abstraction.AccountManager
	Storage        abstraction.Storage
	Production     abstraction.Production
	Market         abstraction.Market
	Tracker        abstraction.Tracker
}

var econ *Economy

func GetEconomyInstance(
	ac abstraction.AccountManager,
	st abstraction.Storage,
	prod abstraction.Production,
	m abstraction.Market,
	tr abstraction.Tracker,
) *Economy {
	if econ != nil {
		return econ
	}
	econ = &Economy{
		AccountManager: ac,
		Storage:        st,
		Production:     prod,
		Market:         m,
		Tracker:        tr,
	}
	return econ
}

func wrap(
	econ *Economy,
	handler func(http.ResponseWriter, *http.Request, *Economy),
) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		handler(w, r, econ)
	}
}

func (econ *Economy) Run() {
	// transfer coins to miners every 10s
	go func() {
		for {
			miners := econ.Storage.GetMiners()
			for _, miner := range miners {
				r := (rand.Float64() * 85) + 30
				mintedCoins := common.DEFAULT_MINT_COINS * r / 100
				econ.AccountManager.PayTo(
					miner.GetAgentID(),
					mintedCoins,
					common.PRIIC,
					common.COIN,
				)
			}
			time.Sleep(time.Duration(common.AVG_BLOCK_MINT_TIME) * time.Second)
		}
	}()

	r := mux.NewRouter()

	r.HandleFunc("/types/{AGENT_TYPE}/agents/join", wrap(econ, Join)).Methods("POST")
	r.HandleFunc("/agents/{AGENT_ID}", wrap(econ, GetAgentByID)).Methods("GET")
	r.HandleFunc("/agents/{AGENT_ID}", wrap(econ, UpdateAgent)).Methods("PUT")

	r.HandleFunc("/agents/{AGENT_ID}/assets", wrap(econ, GetAgentAssets)).Methods("GET")

	// r.HandleFunc("/agents/{AGENT_ID}/wallet/account/balance", wrap(econ, GetWalletAccountBalance)).Methods("GET")
	r.HandleFunc("/agents/{AGENT_ID}/wallet/account", wrap(econ, GetWalletAccount)).Methods("GET")

	r.HandleFunc("/agents/{AGENT_ID}/produce", wrap(econ, Produce)).Methods("POST")

	r.HandleFunc("/agents/{AGENT_ID}/buy", wrap(econ, Buy)).Methods("POST")
	r.HandleFunc("/agents/{AGENT_ID}/sell", wrap(econ, Sell)).Methods("POST")

	r.HandleFunc("/agents/{AGENT_ID}/stabilize", wrap(econ, Stabilize)).Methods("POST")
	r.HandleFunc("/agents/{AGENT_ID}/tokens/buy", wrap(econ, BuyTokens)).Methods("POST")
	r.HandleFunc("/economy/coins/price", wrap(econ, GetCoinPrice)).Methods("GET")
	r.HandleFunc("/economy/tokens/totals", wrap(econ, GetTotalTokens)).Methods("GET")

	fmt.Printf("Listening on port %s", PORT)
	http.ListenAndServe(PORT, r)
}
