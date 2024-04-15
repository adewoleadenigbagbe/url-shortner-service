package enums

import "sort"

type Role int

const (
	Readonly      Role = 1
	User          Role = 2
	Administrator Role = 3
)

var RoleMap = map[Role]string{
	Readonly:      "Readonly",
	User:          "User",
	Administrator: "Administrator",
}

func (p Role) GetValues() []Role {
	var keys []Role
	for k := range RoleMap {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})
	return keys
}

func (r Role) GetKeyValues() []EnumKeyValue[Role, string] {
	var enumKeyValues []EnumKeyValue[Role, string]
	for k, v := range RoleMap {
		enumsKeyValue := EnumKeyValue[Role, string]{
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

func (r Role) String() string {
	return RoleMap[r]
}
