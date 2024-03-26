package helpers

import (
	"bufio"
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
	Cloaking    bool
}

var _ IFileReader = (*TxtReader)(nil)
var _ IFileReader = (*CsvReader)(nil)

type IFileReader interface {
	ReadFile() ([]BulkLinkData, error)
}

type TxtReader struct {
	rc io.ReadCloser
}

func (tReader *TxtReader) ReadFile() ([]BulkLinkData, error) {
	defer tReader.rc.Close()
	scanner := bufio.NewScanner(tReader.rc)

	scanner.Split(bufio.ScanLines)
	var text []string

	for scanner.Scan() {
		text = append(text, scanner.Text())
	}

	datas := lo.Map(text, func(line string, index int) BulkLinkData {
		record := strings.Split(line, ",")
		cloak, _ := strconv.ParseBool(record[2])
		return BulkLinkData{
			OriginalUrl: record[0],
			Alias:       record[1],
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
		cloak, _ := strconv.ParseBool(record[2])
		return BulkLinkData{
			OriginalUrl: record[0],
			Alias:       record[1],
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
		return &TxtReader{
			rc: _rc,
		}, nil
	}

	return nil, errors.New("format not supported")
}
