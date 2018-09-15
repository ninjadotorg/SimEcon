package economy

type Market struct {
	asks []Price
	bids []Price
}

type Price struct {
	price  float32
	orders []Order
}

type Order struct {
	size    float32
	agentId string
	// price   float32
	// filled  float32
	// status  string
}

func (m *Market) trade(side string, size float32, price float32, agentId string) {
}

// markets/{MARKET_ID}/buyLimit?size=&price=&agentId=
func (m *Market) buyLimit(size float32, price float32, agentId string) {
	i, new := findPosition(price, m.bids)
	if new {
		m.bids = insertPosition(i, price, m.bids)
	}
	m.bids[i].orders = append(m.bids[i].orders, Order{size, agentId})
}

// markets/{MARKET_ID}/sellLimit?size=&price=&agentId=
func (m *Market) sellLimit(size float32, price float32, agentId string) {
	i, new := findPosition(price, m.asks)
	if new {
		insertPosition(i, price, m.asks)
	}
	m.asks[i].orders = append(m.asks[i].orders, Order{size, agentId})
}

func findPosition(price float32, prices []Price) (position int, new bool) {
	for i := 0; i < len(prices); i++ {
		if prices[i].price <= price {
			return i, prices[i].price != price
		}
	}
	return len(prices), true
}

func insertPosition(position int, price float32, prices []Price) []Price {
	prices = prices[:len(prices)+1]
	copy(prices[position+1:], prices[position:])
	prices[position] = Price{price: price}
	return prices
}

func (m *Market) buyMarket(size float32, price float32, agentId string) {

}

func (m *Market) sellMarket(size float32, price float32, agentId string) {

}

func (m *Market) bestAsk() float32 {
	return m.asks[0].price
}

func (m *Market) bestBid() float32 {
	return m.bids[0].price
}
