package services

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/adewoleadenigbagbe/url-shortner-service/models"
	"github.com/labstack/echo/v4"
	"github.com/xuri/excelize/v2"
)

var (
	ReportName     = "Links Report"
	ExportFileName = "LinksReport"
)

type ShortData struct {
	OriginalUrl    string
	Hash           string
	Domain         string
	Alias          string
	CreatedOn      time.Time
	ExpirationDate time.Time
	CreatedBy      string
	Cloaking       bool
}

type SheetData struct {
	SheetName  string
	RowCounter int
	Headers    []string
	//TODO: make this generic for the future
	Data []ShortData
}

type ExportService struct {
	Db *sql.DB
}

func (service ExportService) GenerateShortLinkReport(exportContext echo.Context) error {
	var err error
	request := new(models.GenerateShortReportRequest)
	binder := &echo.DefaultBinder{}
	err = binder.BindHeaders(exportContext, request)
	if err != nil {
		return exportContext.JSON(http.StatusBadRequest, []string{err.Error()})
	}

	var organizationName string
	row := service.Db.QueryRow("SELECT Name FROM organizations WHERE Id =? AND IsDeprecated =?", request.OrganizationId, false)
	err = row.Scan(&organizationName)
	if err != nil {
		return exportContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	rows, err := service.Db.Query(`SELECT shortlinks.OriginalUrl,shortlinks.Hash,
	domains.Name,shortlinks.Alias,shortlinks.CreatedOn,shortlinks.ExpirationDate,
	users.Name,shortlinks.Cloaking
	FROM shortlinks 
	JOIN domains on shortlinks.DomainId = domains.Id
	JOIN users on shortlinks.CreatedById = users.Id
	WHERE shortlinks.OrganizationId =? AND shortlinks.IsDeprecated =? AND domains.IsDeprecated =?`,
		request.OrganizationId, false, false)
	if err != nil {
		return exportContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	var shorts []ShortData
	for rows.Next() {
		var short ShortData
		rows.Scan(&short.OriginalUrl, &short.Hash, &short.Domain, &short.Alias, &short.CreatedOn, &short.ExpirationDate, &short.CreatedBy, &short.Cloaking)
		shorts = append(shorts, short)
	}

	defer rows.Close()

	file := excelize.NewFile()
	file.SetActiveSheet(0)
	defer func() {
		if err := file.Close(); err != nil {
			log.Fatal(err)
		}
	}()

	sheetData := CreateSheetData("Links", 1, shorts)
	file.SetSheetName(file.GetSheetName(0), sheetData.SheetName)

	//TODO: make this dynamic , generated from the model in future use
	columnHeaders := []string{
		"OriginalUrl",
		"Short",
		"Domain",
		"Alias",
		"CreatedOn",
		"ExpirationDate",
		"CreatedBy",
		"IsCloak",
	}

	requestedOn := formatSheetDate(time.Now())

	headingInfo := []string{
		organizationName,
		ReportName,
		requestedOn,
	}

	columnLength := calculateSheetWidth(columnHeaders)

	setTitle(file, sheetData, headingInfo, columnLength)

	setColumnHeading(file, sheetData, columnHeaders)

	setDataRows(file, sheetData, columnHeaders)

	buffer, err := file.WriteToBuffer()
	if err != nil {
		return exportContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	attachment := fmt.Sprintf("attachment; filename=%s_%d.xlsx", ExportFileName, time.Now().Nanosecond()*1000)
	exportContext.Response().Header().Set("Content-Disposition", attachment)
	exportContext.Response().Header().Set("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	exportContext.Response().WriteHeader(http.StatusOK)
	exportContext.Response().Write(buffer.Bytes())

	return nil
}

func setTitle(excelFile *excelize.File, sheetData *SheetData, headingInfo []string, columnLength int) {
	var title string
	for _, v := range headingInfo {
		title += fmt.Sprintln(v)
	}

	cell, _ := excelize.CoordinatesToCellName(1, sheetData.RowCounter)
	excelFile.SetCellValue(sheetData.SheetName, cell, title)
	length := len(strings.Split(title, "\n"))
	excelFile.SetRowHeight(sheetData.SheetName, sheetData.RowCounter, float64(length*16))

	style, _ := excelFile.NewStyle(&excelize.Style{
		Font: &excelize.Font{Size: 12, Bold: true},
		Fill: excelize.Fill{Pattern: 1, Color: []string{"##f2f4f5"}, Type: "pattern"},
		Alignment: &excelize.Alignment{
			WrapText: true,
		},
	})

	upperLeftBound, _ := excelize.CoordinatesToCellName(1, sheetData.RowCounter)
	bottomRightBound, _ := excelize.CoordinatesToCellName(columnLength, sheetData.RowCounter)
	excelFile.SetCellStyle(sheetData.SheetName, upperLeftBound, bottomRightBound, style)

	excelFile.MergeCell(sheetData.SheetName, upperLeftBound, bottomRightBound)

	sheetData.NextRow()
}

func calculateSheetWidth(columnHeaders []string) int {
	return len(columnHeaders)
}

func setColumnHeading(excelFile *excelize.File, sheetData *SheetData, columnHeaders []string) {
	style, _ := excelFile.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true},
		Alignment: &excelize.Alignment{Vertical: "center", Horizontal: "center"},
	})
	sheetData.NextRow()
	for columnIndex, columnHeader := range columnHeaders {
		cell, _ := excelize.CoordinatesToCellName(columnIndex+1, sheetData.RowCounter)
		excelFile.SetCellValue(sheetData.SheetName, cell, columnHeader)
		excelFile.SetCellStyle(sheetData.SheetName, cell, cell, style)
	}

	startCol, _ := excelize.ColumnNumberToName(1)
	endCol, _ := excelize.ColumnNumberToName(len(columnHeaders))

	excelFile.SetColWidth(sheetData.SheetName, startCol, endCol, 25)
	sheetData.NextRow()
}

func setDataRows(excelFile *excelize.File, sheetData *SheetData, columnHeaders []string) {
	formatStartRow := sheetData.GetRow()
	fmt.Println(formatStartRow)
	for columnIndex := range columnHeaders {
		for _, d := range sheetData.Data {
			cell, _ := excelize.CoordinatesToCellName(columnIndex+1, sheetData.RowCounter)
			excelFile.SetCellValue(sheetData.SheetName, cell, d)
		}
		sheetData.NextRow()
	}

	//stripped rows
	formatRow()
}

func formatRow() {
}

func (sheetData *SheetData) NextRow() {
	sheetData.RowCounter += 1
}

func (sheetData SheetData) GetRow() int {
	return sheetData.RowCounter
}

func CreateSheetData(name string, rowCounter int, data []ShortData) *SheetData {
	return &SheetData{
		SheetName:  name,
		RowCounter: rowCounter,
		Data:       data,
	}
}

func formatSheetDate(date time.Time) string {
	var (
		month  string
		day    string
		hour   string
		minute string
	)

	if date.Month() < 10 {
		month += "0"
	}
	month += fmt.Sprint(int(date.Month()))

	if date.Day() < 10 {
		day += "0"
	}
	day += fmt.Sprint(date.Day())

	if date.Hour() < 10 {
		hour += "0"
	}
	hour += fmt.Sprint(date.Hour())

	if date.Minute() < 10 {
		minute += "0"
	}
	minute += fmt.Sprint(date.Minute())
	zone, _ := date.Zone()

	return fmt.Sprintf("Requested on: %s/%d/%s %s:%s %s", month, date.Year(), day, hour, minute, zone)
}
