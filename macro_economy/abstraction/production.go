package abstraction

type Production interface {
	GetProductionByAgentType(uint) (AgentProduction, error)
	GetActualAsset(Asset) Asset
	GetActualAssets(map[uint]Asset) map[uint]Asset
}
