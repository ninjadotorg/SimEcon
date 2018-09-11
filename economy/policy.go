package economy

type Policy interface {
	onTick(a *Agent, tick int)
	produce(a *Agent)
	onContract(a *Agent, c Contract)
}

func policy(name string) (p Policy) {
	if name == "household.default" {
		p = &DefaultHousehold{}
	} else if name == "firm.default" {
		p = &DefaultFirm{}
	} else if name == "firm.restaurant" {
		p = &Restaurant{}
	} else if name == "network.contract" {
		p = &ContractExecutor{}
	}
	return
}
