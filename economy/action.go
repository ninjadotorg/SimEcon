package economy

type Action interface {
	init(a *Agent)
	run(a *Agent, s State)
}
