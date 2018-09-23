package storage

type OrderItem struct {
	AgentID      string
	AssetType    uint
	Quantity     float64
	PricePerUnit float64
	OrderTime    int64
}

func (oi *OrderItem) GetOrderTime() int64 {
	return oi.OrderTime
}

func (oi *OrderItem) GetAgentID() string {
	return oi.AgentID
}

func (oi *OrderItem) GetAssetType() uint {
	return oi.AssetType
}

func (oi *OrderItem) GetQuantity() float64 {
	return oi.Quantity
}

func (oi *OrderItem) GetPricePerUnit() float64 {
	return oi.PricePerUnit
}

func (oi *OrderItem) SetQuantity(qty float64) {
	oi.Quantity = qty
}
