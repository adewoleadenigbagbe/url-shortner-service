package helpers

import (
	"encoding/csv"
	"errors"
	"io"
	"strconv"
	"strings"

	"github.com/samber/lo"
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
	rc io.ReadCloser
}

func (tReader *ExcelReader) ReadFile() ([]BulkLinkData, error) {
	defer tReader.rc.Close()
	var text []string
	datas := lo.Map(text, func(line string, index int) BulkLinkData {
		record := strings.Split(line, ",")
		cloak, _ := strconv.ParseBool(record[3])
		return BulkLinkData{
			OriginalUrl: record[0],
			Alias:       record[1],
			Domain:      record[2],
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

	records, err := reader.ReadAll()

	if err != nil {
		return nil, err
	}

	datas := lo.Map(records, func(record []string, index int) BulkLinkData {
		cloak, _ := strconv.ParseBool(record[3])
		return BulkLinkData{
			OriginalUrl: record[0],
			Alias:       record[1],
			Domain:      record[2],
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
