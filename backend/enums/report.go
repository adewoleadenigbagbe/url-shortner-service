package enums

type ReportType string

const (
	Excel ReportType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	Csv   ReportType = "text/csv"
)
