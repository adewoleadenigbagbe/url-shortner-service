package enums

type PayPlan int

const (
	Free PayPlan = iota + 1
	Team
	Enterprise
)
