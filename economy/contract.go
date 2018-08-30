package economy

type Contract struct {
	from         *Agent
	to           *Agent
	payer        *Agent
	payee        *Agent
	amt          float64
	repeat       int
	give         int
	take         int
	status       int // 1: initiated, 2: signed, 3: cancelled
	contractType int // 1: employment
}

func (a *Agent) initiate(c Contract) {
	c.status = 1
	c.to.contract <- c
}

func (a *Agent) handshake(c Contract) {
	c.status = 2
	c.from.contract <- c
}

func (a *Agent) cancel(c Contract) {
	c.status = 3
	c.from.contract <- c
	c.to.contract <- c
}
