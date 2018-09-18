package economy

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/ninjadotorg/SimEcon/util"
)

type Market struct {
	Asks []Price `json:"asks"`
	Bids []Price `json:"bids"`
}

type Price struct {
	Value  float64 `json:"value"`
	Orders []Order `json:"orders"`
}

type Order struct {
	Size    float64 `json:"size"`
	AgentId string  `json:"agentId"`
}

// market/{ASSET_ID}/new
func marketNew(w http.ResponseWriter, r *http.Request) {
	econ.market[mux.Vars(r)["ASSET_ID"]] = &Market{}
}

// market/{ASSET_ID}
func market(w http.ResponseWriter, r *http.Request) {
	if m, ok := econ.market[mux.Vars(r)["ASSET_ID"]]; ok {
		if js, e := json.Marshal(*m); e == nil {
			fmt.Fprintf(w, string(js))
		}
	}
}

// market/{ASSET_ID}/marketBuyLimit?size=&price=&agentId=
func marketBuyLimit(w http.ResponseWriter, r *http.Request) {
	m := econ.market[mux.Vars(r)["ASSET_ID"]]
	q := r.URL.Query()
	if size, e := strconv.ParseFloat(q.Get("size"), 64); e == nil {
		if price, e := strconv.ParseFloat(q.Get("price"), 64); e == nil {
			m.Bids = processLimitOrder("bid", size, price, q.Get("agentId"), m.Bids)
		}
	}
}

// market/{ASSET_ID}/marketSellLimit?size=&price=&agentId=
func marketSellLimit(w http.ResponseWriter, r *http.Request) {
	m := econ.market[mux.Vars(r)["ASSET_ID"]]
	q := r.URL.Query()
	if size, e := strconv.ParseFloat(q.Get("size"), 64); e == nil {
		if price, e := strconv.ParseFloat(q.Get("price"), 64); e == nil {
			m.Asks = processLimitOrder("ask", size, price, q.Get("agentId"), m.Asks)
		}
	}
}

// market/{ASSET_ID}/marketBuy?amount=&agentId=
func marketBuy(w http.ResponseWriter, r *http.Request) {
	assetId := mux.Vars(r)["ASSET_ID"]
	m := econ.market[assetId]
	q := r.URL.Query()
	if amount, e := strconv.ParseFloat(q.Get("amount"), 64); e == nil {
		m.Asks = processOrder("bid", assetId, amount, q.Get("agentId"), m.Asks)
	}
}

// market/{ASSET_ID}/marketSell?amount=&agentId=
func marketSell(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	assetId := mux.Vars(r)["ASSET_ID"]
	agentId := q.Get("agentId")
	m := econ.market[assetId]
	if amount, e := strconv.ParseFloat(q.Get("amount"), 64); e == nil {
		m.Bids = processOrder("ask", assetId, amount, agentId, m.Bids)
	}
}

func processLimitOrder(side string, size float64, price float64, agentId string, prices []Price) []Price {
	for i := 0; i < len(prices); i++ {
		if (side == "bid" && prices[i].Value <= price) || (side == "ask" && prices[i].Value >= price) {
			if prices[i].Value != price {
				prices = append(prices, Price{})
				copy(prices[i+1:], prices[i:])
				prices[i] = Price{Value: price}
			}
			prices[i].Orders = append(prices[i].Orders, Order{size, agentId})
			return prices
		}
	}
	prices = append(prices, Price{Value: price})
	prices[len(prices)-1].Orders = append(prices[len(prices)-1].Orders, Order{size, agentId})
	return prices
}

func processOrder(side string, assetId string, amount float64, agentId string, prices []Price) []Price {
	for i := 0; i < len(prices); i++ {
		price := &prices[i]
		for j := 0; j < len(price.Orders); j++ {
			order := &price.Orders[j]
			if order.Size*price.Value >= amount {
				size := amount / price.Value
				order.Size -= size
				matchOrder(side, assetId, size, price.Value, order.AgentId, agentId)
				j += util.Btoi(order.Size == 0)
				i += util.Btoi(j == len(price.Orders))
				price.Orders = price.Orders[j:]
				prices = prices[i:]
				return prices
			}
			amount -= order.Size * price.Value
			matchOrder(side, assetId, order.Size, price.Value, order.AgentId, agentId)
		}
	}
	if amount > 0 {
		// TODO: maybe return the remaining amount that couldn't buy?
	}
	return prices
}

func matchOrder(side string, assetId string, size float64, price float64, maker string, taker string) {
	if side == "bid" {
		transfer(assetId, size, maker, taker)
		transfer(CASH, size*price, taker, maker)
	} else if side == "ask" {
		transfer(assetId, size, taker, maker)
		transfer(CASH, size*price, maker, taker)
	}
}

func transfer(assetId string, size float64, from string, to string) {
	econ.agent[from].Balance[assetId] -= size
	econ.agent[to].Balance[assetId] += size
}

func (m *Market) bestAsk() float64 {
	return m.Asks[len(m.Asks)-1].Value
}

func (m *Market) bestBid() float64 {
	return m.Bids[0].Value
}
