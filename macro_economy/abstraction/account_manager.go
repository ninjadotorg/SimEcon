package abstraction

type AccountManager interface {
	OpenWalletAccount(string, float64, float64)
	CloseWalletAccount(string)
	GetBalance(string, uint) float64
	GetWalletAccount(string) WalletAccount
	PayFrom(string, float64, uint)
	PayTo(string, float64, int, uint)
	Pay(string, string, float64, int, uint)
	ComputeTotalTokens() map[uint]float64
}
