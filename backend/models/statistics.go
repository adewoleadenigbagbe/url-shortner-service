package models

import (
	"time"

	"github.com/adewoleadenigbagbe/url-shortner-service/enums"
	"github.com/adewoleadenigbagbe/url-shortner-service/helpers/sqltype"
)

type GetShortStatisticRequest struct {
	ShortId        string                      `query:"shortId"`
	OrganizationId string                      `header:"X-OrganizationId"`
	DateRangeType  enums.DateRange             `query:"dateRange"`
	StartDate      sqltype.Nullable[time.Time] `query:"startDate"`
	EndDate        sqltype.Nullable[time.Time] `query:"endDate"`
	Timezone       string                      `query:"timezone"`
}

type GetShortStatisticResponse struct {
}
