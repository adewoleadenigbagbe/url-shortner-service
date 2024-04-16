package enums

import "sort"

type PayPlan int

const (
	Free PayPlan = iota + 1
	Team
	Enterprise
)

var PayPlanMap = map[PayPlan]string{
	Free:       "Free",
	Team:       "Team",
	Enterprise: "Enterprise",
}

func (p PayPlan) GetValues() []PayPlan {
	var keys []PayPlan
	for k := range PayPlanMap {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}

func (p PayPlan) GetKeyValues() []EnumKeyValue[PayPlan, string] {
	var enumKeyValues []EnumKeyValue[PayPlan, string]
	for k, v := range PayPlanMap {
		enumsKeyValue := EnumKeyValue[PayPlan, string]{
			Key:   k,
			Value: v,
		}
		enumKeyValues = append(enumKeyValues, enumsKeyValue)
	}

	sort.Slice(enumKeyValues, func(i, j int) bool {
		return enumKeyValues[i].Key < enumKeyValues[j].Key
	})
	return enumKeyValues
}

func (d PayPlan) String() string {
	return PayPlanMap[d]
}

type PayCycle int

const (
	None PayCycle = iota
	Monthly
	Yearly
)

var PayCycleMap = map[PayCycle]string{
	Monthly: "Monthly",
	Yearly:  "Team",
}

func (d PayCycle) GetValues() []PayCycle {
	var keys []PayCycle
	for k := range PayCycleMap {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}

func (d PayCycle) GetKeyValues() []EnumKeyValue[PayCycle, string] {
	var enumKeyValues []EnumKeyValue[PayCycle, string]
	for k, v := range PayCycleMap {
		enumsKeyValue := EnumKeyValue[PayCycle, string]{
			Key:   k,
			Value: v,
		}
		enumKeyValues = append(enumKeyValues, enumsKeyValue)
	}

	sort.Slice(enumKeyValues, func(i, j int) bool {
		return enumKeyValues[i].Key < enumKeyValues[j].Key
	})
	return enumKeyValues
}

func (d PayCycle) String() string {
	return PayCycleMap[d]
}
