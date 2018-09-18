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

	// SERVER
	PORT = ":8080"
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

func Run() {
	r := mux.NewRouter()

	r.HandleFunc("/agent/{AGENT_ID}", agent)
	r.HandleFunc("/agent/{AGENT_ID}/new", agentNew)
	r.HandleFunc("/agent/{AGENT_ID}/welfare", agentWelfare)
	r.HandleFunc("/agent/{AGENT_ID}/type", agentType)
	r.HandleFunc("/agent/{AGENT_ID}/asset/all", agentAssets)
	r.HandleFunc("/agent/{AGENT_ID}/asset/{ASSET_ID}", agentAsset)
	r.HandleFunc("/agent/{AGENT_ID}/produce", agentProduce)

	r.HandleFunc("/market/{ASSET_ID}", market)
	r.HandleFunc("/market/{ASSET_ID}/new", marketNew)
	r.HandleFunc("/market/{ASSET_ID}/buy", marketBuy)
	r.HandleFunc("/market/{ASSET_ID}/sell", marketSell)
	r.HandleFunc("/market/{ASSET_ID}/buyLimit", marketBuyLimit)
	r.HandleFunc("/market/{ASSET_ID}/sellLimit", marketSellLimit)

	r.HandleFunc("/production/{PRODUCTION_ID}", production)
	r.HandleFunc("/production/{PRODUCTION_ID}/new", productionNew)

	http.ListenAndServe(PORT, r)
}
