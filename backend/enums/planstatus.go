package enums

import "sort"

type PlanStatus int

const (
	Current PlanStatus = iota + 1
	Upcoming
	Archived
)

var PlanStatusMap = map[PlanStatus]string{
	Current:  "Current",
	Upcoming: "Upcoming",
	Archived: "Archived",
}

func (d PlanStatus) GetValues() []PlanStatus {
	var keys []PlanStatus
	for k := range PlanStatusMap {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}

func (d PlanStatus) GetKeyValues() []EnumKeyValue[PlanStatus, string] {
	var enumKeyValues []EnumKeyValue[PlanStatus, string]
	for k, v := range PlanStatusMap {
		enumsKeyValue := EnumKeyValue[PlanStatus, string]{
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

func (d PlanStatus) String() string {
	return PlanStatusMap[d]
}
