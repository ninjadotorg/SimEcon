package dto

type OrderItem struct {
	AgentID      string
	AssetType    uint
	Quantity     float64
	PricePerUnit float64
	OrderTime    int64
}
