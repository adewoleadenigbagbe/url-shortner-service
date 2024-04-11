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
	firstSheetName := "Links"
	f := excelize.NewFile()
	defer f.Close()

	index, _ := f.NewSheet(firstSheetName)
	f.SetActiveSheet(index)
	f.SetSheetName("Sheet1", firstSheetName)

	f.Save()

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

func setHeading() {

}
