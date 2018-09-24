package economy

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
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
}

var econ *Economy

func GetEconomyInstance(
	ac abstraction.AccountManager,
	st abstraction.Storage,
	prod abstraction.Production,
	m abstraction.Market,
) *Economy {
	if econ != nil {
		return econ
	}
	econ = &Economy{
		AccountManager: ac,
		Storage:        st,
		Production:     prod,
		Market:         m,
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
	r := mux.NewRouter()

	r.HandleFunc("/types/{AGENT_TYPE}/agents", wrap(econ, Join)).Methods("POST")
	r.HandleFunc("/agents/{AGENT_ID}/assets", wrap(econ, GetAgentAssets)).Methods("GET")
	r.HandleFunc("/agents/{AGENT_ID}/produce", wrap(econ, Produce)).Methods("POST")
	r.HandleFunc("/agents/{AGENT_ID}/buy", wrap(econ, Buy)).Methods("POST")
	r.HandleFunc("/agents/{AGENT_ID}/sell", wrap(econ, Sell)).Methods("POST")

	fmt.Printf("Listening on port %s", PORT)
	http.ListenAndServe(PORT, r)
}
