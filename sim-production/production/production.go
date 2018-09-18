// Production is a process of combining various material inputs and immaterial
// inputs (plans, know-how) in order to make something for consumption (the
// output). It is the act of creating output, a good or service which has value
// and contributes to the utility of individuals
//
// https://en.wikipedia.org/wiki/Production_(economics)

package production

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

const (
	PORT = ":9090"
)

func Run() {
	r := mux.NewRouter()

	r.HandleFunc("/production/firm", firm)

	http.ListenAndServe(PORT, r)
}

func firm(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var input map[string]float64
	json.Unmarshal([]byte(q.Get("input")), &input)
	// TODO: for now just simply return the input as output
	if js, e := json.Marshal(input); e == nil {
		fmt.Fprintf(w, string(js))
	} else {
		log.Println("err", e)

	}
}
