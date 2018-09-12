package economy

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type Economy struct {
	markets map[string]Market // map[asset]market
	agents  []Agent
}

func (econ *Economy) addAgent(agentType string) (agentId int) {
	agentId = len(econ.agents)
	econ.agents = append(econ.agents, Agent{agentType: agentType})
	return
}

func Run(file string) (e error) {
	r := mux.NewRouter()
	r.HandleFunc("/", homeHandler)
	r.HandleFunc("/agents/add", homeHandler)
	http.ListenAndServe(":8080", r)
	return
}

func newAgent(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello World", vars)

}
func homeHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Hello World", vars)

}
