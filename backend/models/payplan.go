package models

import "github.com/adewoleadenigbagbe/url-shortner-service/enums"

type CreateOrganizationPlanRequest struct {
	PayplanId      string         `json:"payPlanId"`
	PayCycle       enums.PayCycle `json:"payCycle"`
	OrganizationId string         `header:"X-OrganizationId"`
}

type CreateOrganizationPlanResponse struct {
	OrganizationPayPlanId string `json:"organizationPayPlanId"`
}
