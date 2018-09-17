package economy

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Agent struct {
	ProductionId string             `json:"productionId"`
	Balance      map[string]float64 `json:"balance"`
	Welfare      bool               `json:"welfare"`
}

// agent/{AGENT_ID}/new?productionId=
func newAgent(w http.ResponseWriter, r *http.Request) {
	if econ.production[r.URL.Query().Get("productionId")] != nil {
		econ.agent[mux.Vars(r)["AGENT_ID"]] = &Agent{
			ProductionId: r.URL.Query().Get("productionId"),
			Balance:      make(map[string]float64),
		}
	}
}

// agent/{AGENT_ID}/welfare
func welfare(w http.ResponseWriter, r *http.Request) {
	agentId := mux.Vars(r)["AGENT_ID"]
	if !econ.agent[agentId].Welfare {
		if a, ok := econ.agent[agentId]; ok {
			econ.agent[agentId].Welfare = true
			a.Balance[CASH] += WELFARE_AMOUNT
		}
	}
}

// agent/{AGENT_ID}
func agent(w http.ResponseWriter, r *http.Request) {
	if a, ok := econ.agent[mux.Vars(r)["AGENT_ID"]]; ok {
		if js, e := json.Marshal(*a); e == nil {
			fmt.Fprintf(w, string(js))
		}
	}
}

// agent/{AGENT_ID}/type
func agentType(w http.ResponseWriter, r *http.Request) {
	if a, ok := econ.agent[mux.Vars(r)["AGENT_ID"]]; ok {
		fmt.Fprintf(w, a.ProductionId)
	}
}

// agent/{AGENT_ID}/asset/all
func agentAllAssets(w http.ResponseWriter, r *http.Request) {
	if a, ok := econ.agent[mux.Vars(r)["AGENT_ID"]]; ok {
		if js, e := json.Marshal(a.Balance); e == nil {
			fmt.Fprintf(w, string(js))
		}
	}
}

// agent/{AGENT_ID}/asset/{ASSET_ID}
func agentAsset(w http.ResponseWriter, r *http.Request) {
	if a, ok := econ.agent[mux.Vars(r)["AGENT_ID"]]; ok {
		fmt.Fprintf(w, "%f", a.Balance[mux.Vars(r)["ASSET_ID"]])
	}
}

// agent/{AGENT_ID}/produce?input=
func produce(w http.ResponseWriter, r *http.Request) {
	if a, ok := econ.agent[mux.Vars(r)["AGENT_ID"]]; ok {
		p := econ.production[a.ProductionId]
		var input map[string]float64
		json.Unmarshal([]byte(r.URL.Query().Get("input")), &input)
		go p.produce(input, a)
	}
}
