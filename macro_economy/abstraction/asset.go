package abstraction

type Asset interface {
	GetType() uint
	SetQuantity(float64)
	GetProducedTime() int64
	SetProducedTime()
	GetQuantity() float64
}
