// Production is a process of combining various material inputs and immaterial
// inputs (plans, know-how) in order to make something for consumption (the
// output). It is the act of creating output, a good or service which has value
// and contributes to the utility of individuals
//
// https://en.wikipedia.org/wiki/Production_(economics)

package economy

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

type Production struct {
	Type     string             `json:"type"`     // default or custom
	Function string             `json:"function"` // custom endpoint
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

// production/{PRODUCTION_ID}/new?function=&input=&output=&
func productionNew(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	p := Production{Function: q.Get("function")}
	json.Unmarshal([]byte(q.Get("input")), &p.Input)
	json.Unmarshal([]byte(q.Get("output")), &p.Output)
	econ.production[mux.Vars(r)["PRODUCTION_ID"]] = &p
}

func (p *Production) produce(input map[string]float64, agent *Agent) error {

	log.Println("validate input")
	if !sameKeys(input, p.Input) {
		return errors.New("mismatched input")
	}

	log.Println("calculate output")
	output := make(map[string]float64)

	if p.Type == "default" {
		// simply map x to y linearly
		min := math.MaxFloat64
		for k, v := range p.Input {
			r := input[k] / v
			if r < min {
				min = r
			}
		}
		for k, v := range p.Output {
			output[k] = v * min
		}
	} else if p.Type == "custom" {
		// call external endpoint
		url := p.Function
		if js, e := json.Marshal(input); e != nil {
			return e
		} else {
			url += "?input=" + string(js)
		}

		res, _ := http.Get(url)
		data, _ := ioutil.ReadAll(res.Body)

		json.Unmarshal([]byte(string(data)), &output)
	}

	log.Println("validate output")
	if !sameKeys(output, p.Output) {
		return errors.New("mismatched output")
	}

	log.Println("producing")
	time.Sleep(time.Second * time.Duration(p.Time))

	log.Println("update asset balance.. remove input")
	for k, v := range input {
		agent.Balance[k] -= v
	}

	log.Println("update asset balance.. add output")
	for k, v := range output {
		agent.Balance[k] += v
	}

	return nil
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
