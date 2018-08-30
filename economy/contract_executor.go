package economy

type ContractExecutor struct {
	employed int // 0: unemployed, 1: applying, 2: employed
	salary   float64
}

func (d *ContractExecutor) init(a *Agent) {
}

func (d *ContractExecutor) run(a *Agent, s State, econ Economy) {
}

func (d *ContractExecutor) handleContract(a *Agent, c Contract, econ Economy) {
}

func (d *ContractExecutor) checkup(a *Agent, hour int, econ Economy) {
	for _, c := range econ.contracts {
		if hour%c.repeat == 0 {
			c.from.asset.cash -= c.amt
			c.to.asset.cash += c.amt
		}
	}
}
