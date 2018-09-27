package dto

type Person struct {

	// savings rate (portion of total income+savings that is saved in the last step)
	SavingsRate float64

	// consumption (in coin)
	Consumption float64

	// minimum necessity (in real quantity) to buy in the current step
	MinN float64

	// lowest real interest rate seen
	LowRR float64

	// highest real interest rate seen
	HighRR float64

	Agent
}
