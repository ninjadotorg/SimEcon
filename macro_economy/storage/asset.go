package storage

import (
	"time"
)

type Asset struct {
	AgentID      string
	Type         uint
	Quantity     float64
	ProducedTime int64
}

func (a *Asset) GetType() uint {
	return a.Type
}

func (a *Asset) SetQuantity(qty float64) {
	a.Quantity = qty
}

func (a *Asset) GetProducedTime() int64 {
	return a.ProducedTime
}

func (a *Asset) SetProducedTime() {
	a.ProducedTime = time.Now().Unix()
}

func (a *Asset) GetQuantity() float64 {
	return a.Quantity
}
