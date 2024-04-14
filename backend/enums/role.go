package enums

type Role int

const (
	Readonly      Role = 1
	User          Role = 2
	Administrator Role = 3
)

var RoleMap = map[string]Role{
	"Readonly":      Readonly,
	"User":          User,
	"Administrator": Administrator,
}

func (d Role) GetValues() []Role {
	var values []Role
	for _, v := range RoleMap {
		values = append(values, v)
	}

	return values
}

func (p Role) GetValue(key string) Role {
	return RoleMap[key]
}
