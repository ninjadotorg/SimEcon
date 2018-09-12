package economy

import "math"

type Market struct {
	openBuy  []Order
	openSell []Order
}

type Order struct {
	size    float64
	price   float64
	agentId string
	filled  float64
	status  string
}

func (m *Market) trade(side string, size float64, price float64, agentId string) {
	o := Order{}
	o.size = size
	o.price = price
	o.agentId = agentId
	if side == "buy" {
		m.openBuy = append(m.openBuy, o)
	} else if side == "sell" {
		m.openSell = append(m.openSell, o)
	}
}

func (m *Market) buy(size float64, price float64, agentId string) {
	for {
		best := m.bestAsk()
		for _, order := range m.openSell {
			if order.price == best {

			}
		}
	}

}

func (m *Market) bestAsk() (best float64) {
	best = math.MaxFloat64
	for _, order := range m.openSell {
		if order.price < best {
			best = order.price
		}
	}
	return
}

func (m *Market) bestBid() (best float64) {
	best = -math.MaxFloat64
	for _, order := range m.openBuy {
		if order.price > best {
			best = order.price
		}
	}
	return
}
