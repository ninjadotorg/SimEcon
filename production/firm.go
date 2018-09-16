package production

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func Run(file string) {
	r := mux.NewRouter()

	r.HandleFunc("/production/firm", firm)

	http.ListenAndServe(":9090", r)
}

func firm(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var input map[string]float64
	json.Unmarshal([]byte(q.Get("input")), input)
	// TODO: for now just simply return the input as output
	if js, e := json.Marshal(input); e == nil {
		fmt.Fprintf(w, string(js))
	}
}
