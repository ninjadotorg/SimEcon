package economy

type Contract struct {
	from   Agent
	to     Agent
	amt    float64
	repeat int
	give   int
	take   int
}

func (a *Agent) initiate(c Contract) {
	c.to.pendingContracts <- c
}

func (a *Agent) handshake(c Contract) {
	c.from.pendingContracts <- c
}
