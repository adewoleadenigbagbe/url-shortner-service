package enums

type PlanStatus int

const (
	Archived PlanStatus = iota + 1
	Current
	Upcoming
)