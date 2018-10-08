package abstraction

type Agent interface {
	InitAgentAssets(Storage)
	GetAgentID() string
	GetType() uint
	UpdateAgent(Storage, AgentDTO)
}
