package helpers

import (
	"encoding/csv"
	"errors"
	"io"
	"strconv"

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

	rows, err := file.GetRows("Sheet1")
	if err != nil {
		return nil, err
	}

	defer file.Close()

	datas := lo.Map(rows, func(row []string, _ int) BulkLinkData {
		cloak, _ := strconv.ParseBool(row[3])
		return BulkLinkData{
			OriginalUrl: row[0],
			Alias:       row[1],
			Domain:      row[2],
			Cloaking:    cloak,
		}
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

	datas := lo.Map(rows, func(row []string, _ int) BulkLinkData {
		cloak, _ := strconv.ParseBool(row[3])
		return BulkLinkData{
			OriginalUrl: row[0],
			Alias:       row[1],
			Domain:      row[2],
			Cloaking:    cloak,
		}
	})

	return datas, nil
}

func CreateReader(format string, _rc io.ReadCloser) (IFileReader, error) {
	switch format {
	case "text/csv":
		return &CsvReader{
			rc: _rc,
		}, nil
	case "text/plain":
		return &ExcelReader{
			rc: _rc,
		}, nil
	}

	return nil, errors.New("format not supported")
}
