package enums

type ReportType string

const (
	Excel ReportType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	Csv   ReportType = "text/csv"
)

var ReportMap = map[string]ReportType{
	"Excel": Excel,
	"Csv":   Csv,
}

func (d ReportType) GetValues() []ReportType {
	var values []ReportType
	for _, v := range ReportMap {
		values = append(values, v)
	}

	return values
}

func (p ReportType) GetValue(key string) ReportType {
	return ReportMap[key]
}
