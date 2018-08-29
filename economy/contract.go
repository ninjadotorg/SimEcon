package economy

type Contract struct {
	amt    float64
	with   Agent
	repeat int
	give   int
	take   int
}

func (a *Agent) initiate(with Agent, agmt Contract) {

}

func (a *Agent) handshake(with Agent, agmt Contract) {

}
