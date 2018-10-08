package abstraction

type ActionParam interface {
	GetDelta() float64
	GetTax() float64
}
