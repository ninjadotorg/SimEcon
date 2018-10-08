package abstraction

import (
	"sort"

	"github.com/ninjadotorg/SimEcon/common"
)

type OrderItem interface {
	GetOrderTime() int64
	GetAgentID() string
	GetAssetType() uint
	GetQuantity() float64
	GetPricePerUnit() float64
	SetQuantity(float64)
}

type OrderItems []OrderItem

type ByOrderTime struct {
	OrderItems
}

type ByPricePerUnit struct {
	OrderItems
}

func (p OrderItems) Len() int { return len(p) }

func (p OrderItems) Swap(i, j int) { p[i], p[j] = p[j], p[i] }

func (p ByOrderTime) Less(i, j int) bool {
	return p.OrderItems[i].GetOrderTime() < p.OrderItems[j].GetOrderTime()
}

func (p ByPricePerUnit) Less(i, j int) bool {
	return p.OrderItems[i].GetPricePerUnit() < p.OrderItems[j].GetPricePerUnit()
}

func (p OrderItems) SortOrderItems(sortBy uint, isDesc bool) OrderItems {
	if !isDesc {
		if sortBy == common.ORDER_TIME {
			sort.Sort(ByOrderTime{p})
			return p
		}
		if sortBy == common.PRICE_PER_UINT {
			sort.Sort(ByPricePerUnit{p})
			return p
		}
	}
	if sortBy == common.ORDER_TIME {
		sort.Sort(sort.Reverse(ByOrderTime{p}))
		return p
	}
	if sortBy == common.PRICE_PER_UINT {
		sort.Sort(sort.Reverse(ByPricePerUnit{p}))
		return p
	}
	return p
}
