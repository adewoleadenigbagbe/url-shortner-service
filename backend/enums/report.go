package enums

import "sort"

type ReportType string

const (
	Excel ReportType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	Csv   ReportType = "text/csv"
)

var ReportMap = map[ReportType]string{
	Excel: "Excel",
	Csv:   "Csv",
}

func (r ReportType) GetValues() []ReportType {
	var keys []ReportType
	for k := range ReportMap {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}

func (r ReportType) GetKeyValues() []EnumKeyValue[ReportType, string] {
	var enumKeyValues []EnumKeyValue[ReportType, string]
	for k, v := range ReportMap {
		enumsKeyValue := EnumKeyValue[ReportType, string]{
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

func (r ReportType) String() string {
	return ReportMap[r]
}
