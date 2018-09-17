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
	market     map[string]*Market     // map[asset_id]market
	agent      map[string]*Agent      // map[agent_id]agent
	production map[string]*Production // map[agent_id]agent
}

var econ Economy = Economy{
	market:     make(map[string]*Market),
	agent:      make(map[string]*Agent),
	production: make(map[string]*Production),
}

func Run(file string) {
	r := mux.NewRouter()

	r.HandleFunc("/agent/{AGENT_ID}", agent)
	r.HandleFunc("/agent/{AGENT_ID}/new", newAgent)
	r.HandleFunc("/agent/{AGENT_ID}/welfare", welfare)
	r.HandleFunc("/agent/{AGENT_ID}/type", agentType)
	r.HandleFunc("/agent/{AGENT_ID}/asset/all", agentAllAssets)
	r.HandleFunc("/agent/{AGENT_ID}/asset/{ASSET_ID}", agentAsset)

	r.HandleFunc("/market/{ASSET_ID}", market)
	r.HandleFunc("/market/{ASSET_ID}/new", newMarket)
	r.HandleFunc("/market/{ASSET_ID}/buy", buy)
	r.HandleFunc("/market/{ASSET_ID}/sell", sell)
	r.HandleFunc("/market/{ASSET_ID}/buyLimit", buyLimit)
	r.HandleFunc("/market/{ASSET_ID}/sellLimit", sellLimit)

	r.HandleFunc("/production/{PRODUCTION_ID}", production)
	r.HandleFunc("/production/{PRODUCTION_ID}/new", newProduction)

	http.ListenAndServe(":8080", r)
}
