package account_manager

import (
	"sync"

	"github.com/ninjadotorg/SimEcon/common"
	"github.com/ninjadotorg/SimEcon/macro_economy/abstraction"
)

var accountManager *AccountManager

type AccountManager struct {
	locker         *sync.RWMutex
	WalletAccounts map[string]*WalletAccount
	// TODO: bankbook
}

func GetAccountManagerInstance() *AccountManager {
	if accountManager != nil {
		return accountManager
	}
	accountManager = &AccountManager{
		locker:         &sync.RWMutex{},
		WalletAccounts: map[string]*WalletAccount{},
	}
	return accountManager
}

func (accManager *AccountManager) OpenWalletAccount(
	agentID string,
	coins float64,
	bonds float64,
) {
	accManager.locker.Lock()
	defer accManager.locker.Unlock()
	newAddress := common.UUID()
	acc := NewWalletAccount(newAddress, coins, bonds)
	accManager.WalletAccounts[agentID] = acc
}

func (accManager *AccountManager) CloseWalletAccount(
	agentID string,
) {
	accManager.locker.Lock()
	defer accManager.locker.Unlock()
	if _, ok := accManager.WalletAccounts[agentID]; ok {
		delete(accManager.WalletAccounts, agentID)
	}
}

func (accManager *AccountManager) GetBalance(
	agentID string,
	tokenType uint,
) float64 {
	accManager.locker.RLock()
	defer accManager.locker.RUnlock()
	acc, ok := accManager.WalletAccounts[agentID]
	if !ok {
		return 0
	}
	if tokenType == common.COIN {
		return acc.Coins
	}
	if tokenType == common.BOND {
		return acc.Bonds
	}
	return 0
}

func (accManager *AccountManager) GetWalletAccount(
	agentID string,
) abstraction.WalletAccount {
	accManager.locker.RLock()
	defer accManager.locker.RUnlock()
	acc, ok := accManager.WalletAccounts[agentID]
	if !ok {
		return nil
	}
	return acc
}

func (accManager *AccountManager) PayFrom(
	payerID string,
	amt float64,
	tokenType uint,
) {
	accManager.locker.Lock()
	defer accManager.locker.Unlock()
	fromAcc := accManager.WalletAccounts[payerID]
	// if fromAcc.Balance < amt { // TODO: will handle this case later

	// }
	if tokenType == common.COIN {
		fromAcc.Coins -= amt
		return
	}
	if tokenType == common.BOND {
		fromAcc.Bonds -= amt
	}
}

func (accManager *AccountManager) PayTo(
	payeeID string,
	amt float64,
	purpose int, // either PRIIC or SECIC
	tokenType uint,
) {
	accManager.locker.Lock()
	defer accManager.locker.Unlock()
	toAcc := accManager.WalletAccounts[payeeID]
	if tokenType != common.COIN && tokenType != common.BOND {
		return
	}

	if tokenType == common.COIN {
		toAcc.Coins += amt
	} else {
		toAcc.Bonds += amt
	}
	if purpose == common.PRIIC {
		toAcc.PriIC = amt
		return
	}
	toAcc.SecIC = amt
}

func (accManager *AccountManager) Pay(
	payerID string,
	payeeID string,
	amt float64,
	purpose int, // either PRIIC or SECIC
	tokenType uint,
) {
	accManager.locker.Lock()
	defer accManager.locker.Unlock()
	fromAcc := accManager.WalletAccounts[payerID]
	toAcc := accManager.WalletAccounts[payeeID]
	if tokenType != common.COIN && tokenType != common.BOND {
		return
	}

	if tokenType == common.COIN {
		fromAcc.Coins -= amt
		toAcc.Coins += amt
	} else {
		fromAcc.Bonds -= amt
		toAcc.Bonds += amt
	}

	if purpose == common.PRIIC {
		toAcc.PriIC = amt
		return
	}
	toAcc.SecIC = amt
}

func (am *AccountManager) ComputeTotalTokens() map[uint]float64 {
	totalTokens := map[uint]float64{}
	var totalCoins float64 = 0
	var totalBonds float64 = 0
	for _, wlAcc := range am.WalletAccounts {
		totalCoins += wlAcc.Coins
		totalBonds += wlAcc.Bonds
	}
	totalTokens[common.COIN] = totalCoins
	totalTokens[common.BOND] = totalBonds
	return totalTokens
}
