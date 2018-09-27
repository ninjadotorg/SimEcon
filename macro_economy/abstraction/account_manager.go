package abstraction

type AccountManager interface {
	OpenWalletAccount(string, float64)
	CloseWalletAccount(string)
	GetBalance(string) float64
	GetWalletAccount(string) WalletAccount
	Pay(string, string, float64, int)
}
