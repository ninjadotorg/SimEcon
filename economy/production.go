// Production is a process of combining various material inputs and immaterial
// inputs (plans, know-how) in order to make something for consumption (the
// output). It is the act of creating output, a good or service which has value
// and contributes to the utility of individuals
//
// https://en.wikipedia.org/wiki/Production_(economics)

package economy

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Production struct {
	Function string             `json:"function"` // end point to function
	Time     int                `json:"time"`     // how long it takes to produce
	Input    map[string]float64 `json:"input"`    // []map(asset_id)size
	Output   map[string]float64 `json:"output"`   // []map(asset_id)size
}

// production/{PRODUCTION_ID}
func production(w http.ResponseWriter, r *http.Request) {
	if r, ok := econ.production[mux.Vars(r)["PRODUCTION_ID"]]; ok {
		if js, e := json.Marshal(r); e == nil {
			fmt.Fprintf(w, string(js))
		}
	}
}

// production/add?function=&input=&output=&
func addProduction(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	p := Production{Function: q.Get("function")}
	json.Unmarshal([]byte(q.Get("input")), &p.Input)
	json.Unmarshal([]byte(q.Get("output")), &p.Output)
	econ.production[mux.Vars(r)["PRODUCTION_ID"]] = &p
}

func (p *Production) produce(input map[string]float64, agent *Agent) {
	// validate input
	if !sameKeys(input, p.Input) {
		return
	}

	// calc output
	log.Println("Calling endpoint", p.Function, "with ", input)
	var output map[string]float64

	// validate output
	if !sameKeys(output, p.Output) {
		return
	}

	// producing...
	time.Sleep(time.Second * time.Duration(p.Time))

	// minus input, add output
	for k, v := range input {
		agent.Balance[k] -= v
	}
	for k, v := range output {
		agent.Balance[k] += v
	}
}

func sameKeys(a, b map[string]float64) bool {
	if len(a) != len(b) {
		return false
	}
	for k, _ := range a {
		if _, ok := b[k]; !ok {
			return false
		}
	}
	return true
}
