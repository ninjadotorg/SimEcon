package abstraction

type Production interface {
	GetProductionByAgentType(uint) (AgentProduction, error)
}
