package economy

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Agent struct {
	Type    string             `json:"type"`
	Balance map[string]float64 `json:"balance"`
}

// agents/{AGENT_ID}/join?type=
func join(w http.ResponseWriter, r *http.Request) {
	econ.agent[mux.Vars(r)["AGENT_ID"]] = &Agent{
		Type:    r.URL.Query().Get("type"),
		Balance: make(map[string]float64),
	}
}

// agents/{AGENT_ID}/welfare
func welfare(w http.ResponseWriter, r *http.Request) {
	agentId := mux.Vars(r)["AGENT_ID"]
	if !econ.welfare[agentId] {
		if a, ok := econ.agent[agentId]; ok {
			econ.welfare[agentId] = true
			a.Balance[CASH] += WELFARE_AMOUNT
		}
	}
}

// agents/{AGENT_ID}
func agent(w http.ResponseWriter, r *http.Request) {
	if a, ok := econ.agent[mux.Vars(r)["AGENT_ID"]]; ok {
		if js, e := json.Marshal(a); e == nil {
			fmt.Fprintf(w, string(js))
		}
	}
}

// agents/{AGENT_ID}/type
func agentType(w http.ResponseWriter, r *http.Request) {
	if a, ok := econ.agent[mux.Vars(r)["AGENT_ID"]]; ok {
		fmt.Fprintf(w, a.Type)
	}
}

// agents/{AGENT_ID}/assets
func agentAssets(w http.ResponseWriter, r *http.Request) {
	if a, ok := econ.agent[mux.Vars(r)["AGENT_ID"]]; ok {
		if js, e := json.Marshal(a.Balance); e == nil {
			fmt.Fprintf(w, string(js))
		}
	}
}

// agents/{AGENT_ID}/assets/{ASSET_ID}
func agentAssetBalance(w http.ResponseWriter, r *http.Request) {
	if a, ok := econ.agent[mux.Vars(r)["AGENT_ID"]]; ok {
		fmt.Fprintf(w, "%f", a.Balance[mux.Vars(r)["ASSET_ID"]])
	}
}
