package account_manager

type WalletAccount struct {
	Address string
	Balance float64 // checking/saving acc balance
	PriIC   float64 // primary income in the last step
	SecIC   float64 // secondary income in the last step
}

func NewWalletAccount(address string, balance float64) *WalletAccount {
	return &WalletAccount{
		Address: address,
		Balance: balance,
		PriIC:   0,
		SecIC:   0,
	}
}
