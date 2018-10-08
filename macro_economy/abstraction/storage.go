package abstraction

type Storage interface {
	InsertAgent(string, uint) Agent
	UpdateAssets(string, []Asset)
	UpdateAsset(string, Asset)
	GetAgentAssets(string) (map[uint]Asset, error)
	GetAgentAsset(string, uint) (Asset, error)
	GetAgentByID(string) (Agent, error)
	GetSortedBidsByAssetType(uint, bool, uint) OrderItems
	RemoveBidsByAgentIDs([]string, uint) error
	AppendAsk(uint, string, float64, float64)
	GetSortedAsksByAssetType(uint, bool, uint) OrderItems
	RemoveAsksByAgentIDs([]string, uint) error
	AppendBid(uint, string, float64, float64)
	GetTotalAsksByAssetType(uint) float64
	GetTotalBidsByAssetType(uint) float64
	ComputeDecidedParam() ActionParam
	InsertParam(string, float64, float64) ActionParam
	GetMiners() []Agent
}
