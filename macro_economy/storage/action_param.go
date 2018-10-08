package storage

type ActionParam struct {
	AgentID      string
	Delta        float64
	Tax          float64
	ModifiedTime int64
}

func (ap *ActionParam) GetDelta() float64 {
	return ap.Delta
}

func (ap *ActionParam) GetTax() float64 {
	return ap.Tax
}
