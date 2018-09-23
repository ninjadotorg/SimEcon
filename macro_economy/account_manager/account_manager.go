package account_manager

import (
	"github.com/ninjadotorg/SimEcon/common"
)

var accountManager *AccountManager

type AccountManager struct {
	WalletAccounts map[string]*WalletAccount
	// TODO: bankbook
}

func GetAccountManagerInstance() *AccountManager {
	if accountManager != nil {
		return accountManager
	}
	accountManager = &AccountManager{
		WalletAccounts: map[string]*WalletAccount{},
	}
	return accountManager
}

func (accManager *AccountManager) OpenWalletAccount(
	agentID string,
	balance float64,
) {
	newAddress := common.UUID()
	acc := NewWalletAccount(newAddress, balance)
	accManager.WalletAccounts[agentID] = acc
}

func (accManager *AccountManager) CloseWalletAccount(
	agentID string,
) {
	if _, ok := accManager.WalletAccounts[agentID]; ok {
		delete(accManager.WalletAccounts, agentID)
	}
}

func (accManager *AccountManager) GetBalance(
	agentID string,
) float64 {
	acc, ok := accManager.WalletAccounts[agentID]
	if !ok {
		return 0
	}
	return acc.Balance
}

func (accManager *AccountManager) PayFrom(
	payerID string,
	amt float64,
) {
	fromAcc := accManager.WalletAccounts[payerID]
	// if fromAcc.Balance < amt { // TODO: will handle this case later

	// }
	fromAcc.Balance -= amt
}

func (accManager *AccountManager) PayTo(
	payeeID string,
	amt float64,
	purpose int, // either PRIIC or SECIC
) {
	toAcc := accManager.WalletAccounts[payeeID]
	toAcc.Balance += amt
	if purpose == common.PRIIC {
		toAcc.PriIC += amt
		return
	}
	toAcc.SecIC += amt
}

func (accManager *AccountManager) Pay(
	payerID string,
	payeeID string,
	amt float64,
	purpose int, // either PRIIC or SECIC
) {
	fromAcc := accManager.WalletAccounts[payerID]
	toAcc := accManager.WalletAccounts[payeeID]

	fromAcc.Balance -= amt
	toAcc.Balance += amt
	if purpose == common.PRIIC {
		toAcc.PriIC += amt
		return
	}
	toAcc.SecIC += amt
}
