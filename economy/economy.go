package economy

import (
	"net/http"

	"github.com/gorilla/mux"
)

const (
	// ASSETS
	CASH = "CASH"

	// ENV VARS
	WELFARE_AMOUNT = 10
)

type Economy struct {
	market  map[string]*Market // map[asset_id]market
	agent   map[string]*Agent  // map[agent_id]agent
	welfare map[string]bool
}

var econ Economy = Economy{
	market:  make(map[string]*Market),
	agent:   make(map[string]*Agent),
	welfare: make(map[string]bool),
}

func Run(file string) {
	r := mux.NewRouter()

	r.HandleFunc("/agents/{AGENT_ID}", agent)
	r.HandleFunc("/agents/{AGENT_ID}/join", join)
	r.HandleFunc("/agents/{AGENT_ID}/welfare", welfare)
	r.HandleFunc("/agents/{AGENT_ID}/type", agentType)
	r.HandleFunc("/agents/{AGENT_ID}/assets", agentAssets)
	r.HandleFunc("/agents/{AGENT_ID}/assets/{ASSET_ID}", agentAssetBalance)

	r.HandleFunc("/markets/{ASSET_ID}/buy", buy)
	r.HandleFunc("/markets/{ASSET_ID}/sell", sell)
	r.HandleFunc("/markets/{ASSET_ID}/buyLimit", buyLimit)
	r.HandleFunc("/markets/{ASSET_ID}/sellLimit", sellLimit)

	http.ListenAndServe(":8080", r)
}
