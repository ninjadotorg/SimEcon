package economy

import (
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ninjadotorg/SimEcon/util"
)

type Market struct {
	asks []Price
	bids []Price
}

type Price struct {
	value  float64
	orders []Order
}

type Order struct {
	size    float64
	agentId string
}

// market/{ASSET_ID}/buyLimit?size=&price=&agentId=
func buyLimit(w http.ResponseWriter, r *http.Request) {
	m := econ.market[mux.Vars(r)["ASSET_ID"]]
	q := r.URL.Query()
	if size, e := strconv.ParseFloat(q.Get("size"), 64); e == nil {
		if price, e := strconv.ParseFloat(q.Get("price"), 64); e == nil {
			m.bids = processLimitOrder(size, price, q.Get("agentId"), m.bids)
		}
	}
}

// market/{ASSET_ID}/sellLimit?size=&price=&agentId=
func sellLimit(w http.ResponseWriter, r *http.Request) {
	m := econ.market[mux.Vars(r)["ASSET_ID"]]
	q := r.URL.Query()
	if size, e := strconv.ParseFloat(q.Get("size"), 64); e == nil {
		if price, e := strconv.ParseFloat(q.Get("price"), 64); e == nil {
			m.asks = processLimitOrder(size, price, q.Get("agentId"), m.asks)
		}
	}
}

// market/{ASSET_ID}/buy?amount=&agentId=
func buy(w http.ResponseWriter, r *http.Request) {
	m := econ.market[mux.Vars(r)["ASSET_ID"]]
	q := r.URL.Query()
	if amount, e := strconv.ParseFloat(q.Get("amount"), 64); e == nil {
		m.asks = processOrder(amount, q.Get("agentId"), m.asks)
	}
}

// market/{ASSET_ID}/sell?amount=&agentId=
func sell(w http.ResponseWriter, r *http.Request) {
	m := econ.market[mux.Vars(r)["ASSET_ID"]]
	q := r.URL.Query()
	if amount, e := strconv.ParseFloat(q.Get("amount"), 64); e == nil {
		m.bids = processOrder(amount, q.Get("agentId"), m.bids)
	}
}

func processLimitOrder(size float64, price float64, agentId string, prices []Price) []Price {
	for i := 0; i < len(prices); i++ {
		if prices[i].value <= price {
			if prices[i].value == price {
				prices = prices[:len(prices)+1]
				copy(prices[i+1:], prices[i:])
				prices[i] = Price{value: price}
			}
			prices[i].orders = append(prices[i].orders, Order{size, agentId})
			return prices
		}
	}
	return prices
}

func processOrder(amount float64, agentId string, prices []Price) []Price {
	for i := 0; i < len(prices); i++ {
		price := prices[i]
		for j := 0; j < len(price.orders); j++ {
			order := price.orders[j]
			if order.size*price.value >= amount {
				order.size -= amount / price.value
				j += util.Btoi(order.size == 0)
				i += util.Btoi(j == len(price.orders))
				price.orders = price.orders[j:]
				prices = prices[i:]
				return prices
			}
			amount -= order.size * price.value
		}
	}
	if amount > 0 {
		// TODO: maybe return the remaining amount that couldn't buy?
	}
	return prices
}

func (m *Market) bestAsk() float64 {
	return m.asks[len(m.asks)-1].value
}

func (m *Market) bestBid() float64 {
	return m.bids[0].value
}
