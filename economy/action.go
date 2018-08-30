package economy

type Action interface {
	init(a *Agent)
	run(a *Agent, s State, econ *Economy)
	handleContract(a *Agent, c Contract, econ *Economy)
	checkup(a *Agent, hour int, econ *Economy)
}
