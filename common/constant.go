package common

const (
	PRIIC = 1
	SECIC = 2

	// Agent types
	PERSON          = 1
	NECESSITY_FIRM  = 2
	CAPITAL_FIRM    = 3
	COMMERCIAL_BANK = 4

	// Agent asset types
	NECESSITY = 1
	CAPITAL   = 2
	MAN_HOUR  = 3

	// Default initial assets
	PERSON_NECESSITY = 10
	PERSON_MAN_HOURS = 0

	NFIRM_MAN_HOURS = 10
	NFIRM_CAPITAL   = 10
	NFIRM_NECESSITY = 0

	CAPITAL_MAN_HOURS = 10
	CAPITAL_CAPITAL   = 0

	// decay params
	NECESSITY_DECAY_PERIOD  = 300 // 5M
	NECESSITY_EPSILON_DECAY = 0.95
	CAPITAL_DECAY_PERIOD    = 240 // 4M
	CAPITAL_EPSILON_DECAY   = 0.975
	MAN_HOURS_DECAY_PERIOD  = 360 // 6M
	MAN_HOURS_EPSILON_DECAY = 0.9
)