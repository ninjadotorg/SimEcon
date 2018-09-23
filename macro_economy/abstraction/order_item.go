package abstraction

import "sort"

type OrderItem interface {
	GetOrderTime() int64
	GetAgentID() string
	GetAssetType() uint
	GetQuantity() float64
	GetPricePerUnit() float64
	SetQuantity(float64)
}

type OrderItems []OrderItem

func (p OrderItems) Len() int           { return len(p) }
func (p OrderItems) Less(i, j int) bool { return p[i].GetOrderTime() < p[j].GetOrderTime() }
func (p OrderItems) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

func (p OrderItems) SortOrderItems(isDesc bool) OrderItems {
	if isDesc {
		sort.Sort(sort.Reverse(p))
	}
	sort.Sort(p)
	return p
}
