package helpers

import (
	"encoding/csv"
	"errors"
	"io"
	"strconv"

	"github.com/adewoleadenigbagbe/url-shortner-service/enums"
	"github.com/samber/lo"
	"github.com/xuri/excelize/v2"
)

type BulkLinkData struct {
	OriginalUrl string
	Alias       string
	Domain      string
	Cloaking    bool
}

var _ IFileReader = (*ExcelReader)(nil)
var _ IFileReader = (*CsvReader)(nil)

type IFileReader interface {
	ReadFile() ([]BulkLinkData, error)
}

type ExcelReader struct {
	rc io.Reader
}

func (excelReader *ExcelReader) ReadFile() ([]BulkLinkData, error) {
	file, err := excelize.OpenReader(excelReader.rc)
	if err != nil {
		return nil, err
	}

	rows, err := file.GetRows(file.GetSheetName(0))
	if err != nil {
		return nil, err
	}

	defer file.Close()
	datas := lo.FilterMap(rows, func(row []string, index int) (BulkLinkData, bool) {
		if index == 0 {
			return BulkLinkData{}, false
		}

		cloak, _ := strconv.ParseBool(row[3])
		return BulkLinkData{
			OriginalUrl: row[0],
			Alias:       row[1],
			Domain:      row[2],
			Cloaking:    cloak,
		}, true
	})

	return datas, nil
}

type CsvReader struct {
	rc io.ReadCloser
}

func (csvReader *CsvReader) ReadFile() ([]BulkLinkData, error) {
	defer csvReader.rc.Close()
	reader := csv.NewReader(csvReader.rc)

	rows, err := reader.ReadAll()

	if err != nil {
		return nil, err
	}

	datas := lo.FilterMap(rows, func(row []string, index int) (BulkLinkData, bool) {
		if index == 0 {
			return BulkLinkData{}, false
		}

		cloak, _ := strconv.ParseBool(row[3])
		return BulkLinkData{
			OriginalUrl: row[0],
			Alias:       row[1],
			Domain:      row[2],
			Cloaking:    cloak,
		}, true
	})

	return datas, nil
}

func CreateReader(format enums.ReportType, _rc io.ReadCloser) (IFileReader, error) {
	switch format {
	case enums.Csv:
		return &CsvReader{
			rc: _rc,
		}, nil
	case enums.Excel:
		return &ExcelReader{
			rc: _rc,
		}, nil
	}

	return nil, errors.New("format not supported")
}
