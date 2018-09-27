package abstraction

type Agent interface {
	InitAgentAssets(Storage)
	GetType() uint
	UpdateAgent(Storage, AgentDTO)
}
