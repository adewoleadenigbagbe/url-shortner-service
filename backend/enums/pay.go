package enums

type PayType int

const (
	Premium PayType = iota + 1
	Standard
	Gold
)
