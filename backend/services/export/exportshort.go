package services

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/xuri/excelize/v2"
)

var (
	ReportName     = "Links Report"
	ExportFileName = "LinksReport"
)

type ExportService struct {
	Db *sql.DB
}

func (service ExportService) GenerateShortLinkReport(exportContext echo.Context) error {
	f := excelize.NewFile()
	defer f.Close()

	sheetData := CreateSheet("Sheet1", "Links", 0)
	index, _ := f.NewSheet(sheetData.SheetId)
	f.SetActiveSheet(index)
	f.SetSheetName(sheetData.SheetId, sheetData.SheetName)

	buffer, err := f.WriteToBuffer()
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

func setHeading(excelFile *excelize.File, sheetData *SheetData) error {
	return nil
}

type SheetData struct {
	SheetId    string
	SheetName  string
	RowCounter *int
	Headers    []string
	Data       [][]interface{}
}

func CreateSheet(id, name string, rowCounter int) *SheetData {
	return &SheetData{
		SheetId:    id,
		SheetName:  name,
		RowCounter: &rowCounter,
		Data:       make([][]interface{}, 0),
	}
}
