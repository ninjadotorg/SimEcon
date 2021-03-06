package common

const (
	DEFAULT_AGENT_ID = "00000000-00000000-00000000-00000000"

	PRIIC = 1
	SECIC = 2

	// Agent types
	PERSON          = 1
	NECESSITY_FIRM  = 2
	CAPITAL_FIRM    = 3
	COMMERCIAL_BANK = 4
	PRICE_STABILITY = 5
	MINER           = 6

	// Miner status
	CANDIDATE = 1
	LEADER    = 2

	// Agent asset types
	NECESSITY = 1
	CAPITAL   = 2
	MAN_HOUR  = 3
	COIN      = 4
	BOND      = 5

	// Default initial assets
	PERSON_NECESSITY = 8
	PERSON_MAN_HOURS = 0

	NFIRM_MAN_HOURS = 0
	NFIRM_CAPITAL   = 0
	NFIRM_NECESSITY = 0

	CAPITAL_MAN_HOURS = 40
	CAPITAL_CAPITAL   = 0

	MINER_NECESSITY = 8

	// decay params
	NECESSITY_DECAY_PERIOD  = 600 //300 // 5M
	NECESSITY_EPSILON_DECAY = 0.95
	CAPITAL_DECAY_PERIOD    = 600 //240 // 4M
	CAPITAL_EPSILON_DECAY   = 0.975
	MAN_HOURS_DECAY_PERIOD  = 600 //360 // 6M
	MAN_HOURS_EPSILON_DECAY = 0.9

	// init account balance
	DEFAULT_ACCOUNT_BALANCE = 0
	NFIRM_ACCOUNT_BALANCE   = 800

	// baseline of assets' price
	// NECESSITY_PRICE_BASELINE = 5
	// CAPITAL_PRICE_BASELINE   = 12
	// MAN_HOUR_PRICE_BASELINE  = 8

	NECESSITY_PRICE_BASELINE = 8
	CAPITAL_PRICE_BASELINE   = 8
	MAN_HOUR_PRICE_BASELINE  = 8

	// order item sort fields
	PRICE_PER_UINT = 1
	ORDER_TIME     = 2

	// Default mint coins
	DEFAULT_MINT_COINS  = 10
	AVG_BLOCK_MINT_TIME = 10

	// tracking path
	TOTAL_ASKS_FILE = "/Users/autonomous/projects/golang-projects/src/github.com/ninjadotorg/SimEcon/macro_economy/market/total_asks"
	TOTAL_BIDS_FILE = "/Users/autonomous/projects/golang-projects/src/github.com/ninjadotorg/SimEcon/macro_economy/market/total_bids"
)
