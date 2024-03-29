package models

const (
	ApiExpiry        = 1
	EmailRegex       = "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9-]+(?:\\.[a-zA-Z0-9-]+)*$"
	PhoneNumberRegex = "\\+[1-9]{1}[0-9]{0,2}-[2-9]{1}[0-9]{2}-[2-9]{1}[0-9]{2}-[0-9]{4}$"

	//users
	Team_Plan_User_Limit = 1000

	//Domain
	Free_Plan_Domain_Limit = 5
	Team_Plan_Domain_Limit = 50
)
