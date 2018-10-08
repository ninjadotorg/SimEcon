package account_manager

type WalletAccount struct {
	Address string
	Coins   float64
	Bonds   float64
	PriIC   float64 // primary income in the last step
	SecIC   float64 // secondary income in the last step
}

func NewWalletAccount(
	address string,
	coins float64,
	bonds float64,
) *WalletAccount {
	return &WalletAccount{
		Address: address,
		Coins:   coins,
		Bonds:   bonds,
		PriIC:   0,
		SecIC:   0,
	}
}
