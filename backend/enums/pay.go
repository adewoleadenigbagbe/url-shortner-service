package enums

type PayPlan int

const (
	Free PayPlan = iota + 1
	Team
	Enterprise
)

var PayPlanMap = map[string]PayPlan{
	"Free":       Free,
	"Team":       Team,
	"Enterprise": Enterprise,
}

func (d PayPlan) GetValues() []PayPlan {
	var values []PayPlan
	for _, v := range PayPlanMap {
		values = append(values, v)
	}

	return values
}

func (p PayPlan) GetValue(key string) PayPlan {
	return PayPlanMap[key]
}

type PayCycle int

const (
	None PayCycle = iota
	Monthly
	Yearly
)

var PayCycleMap = map[string]PayCycle{
	"Monthly": Monthly,
	"Team":    Yearly,
}

func (d PayCycle) GetValues() []PayCycle {
	var values []PayCycle
	for _, v := range PayCycleMap {
		values = append(values, v)
	}

	return values
}

func (p PayCycle) GetValue(key string) PayCycle {
	return PayCycleMap[key]
}
