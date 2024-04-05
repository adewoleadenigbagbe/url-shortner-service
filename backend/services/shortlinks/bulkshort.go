package services

import (
	"math"
	"net/http"
	"time"

	sequentialguid "github.com/adewoleadenigbagbe/sequential-guid"
	"github.com/adewoleadenigbagbe/url-shortner-service/helpers"
	"github.com/labstack/echo/v4"
)

func (service UrlService) CreateBulkShortLink(urlContext echo.Context) error {
	var err error
	req := urlContext.Request()
	headers := req.Header

	reader, err := helpers.CreateReader(headers.Get("Content-Type"), req.Body)
	if err != nil {
		return urlContext.JSON(http.StatusBadRequest, []string{err.Error()})
	}

	data, err := reader.ReadFile()
	if err != nil {
		return urlContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	organizationId := headers.Get("X-OrganizationId")
	userId := headers.Get("X-UserId")

	rows, err := service.Db.Query("SELECT Id, Name FROM domains WHERE OrganizationId =? AND IsDeprecated =?", organizationId, false)
	if err != nil {
		return urlContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	defer rows.Close()

	domainMap := make(map[string]string)
	for rows.Next() {
		var id string
		var name string
		err = rows.Scan(&id, &name)
		if err != nil {
			return urlContext.JSON(http.StatusInternalServerError, []string{err.Error()})
		}
		domainMap[name] = id
	}

	var batchSize int = 500
	batchNumber := 0
	total := len(data)
	remaining := total
	totalBatch := int(math.Ceil(float64(total) / float64(batchSize)))

	tx, err := service.Db.Begin()
	if err != nil {
		return urlContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	for batchNumber < totalBatch {
		if remaining < batchSize {
			batchSize = remaining
		}
		stmt, args := insertBulkShortStmt(organizationId, userId, data[batchNumber*batchSize:(batchNumber*batchSize)+batchSize], domainMap)
		_, err = tx.Exec(stmt, args...)
		if err != nil {
			tx.Rollback()
		}
		batchNumber += 1
		remaining -= batchSize
	}

	tx.Commit()

	return urlContext.JSON(http.StatusCreated, nil)
}

func insertBulkShortStmt(organizationId string, userId string, data []helpers.BulkLinkData, dic map[string]string) (string, []interface{}) {
	stmt := "INSERT INTO shortlinks VALUES "
	vals := []interface{}{}
	for _, d := range data {
		id := sequentialguid.NewSequentialGuid().String()
		short := helpers.GenerateShortLink(d.OriginalUrl)
		now := time.Now()
		expirationDate := now.AddDate(expirySpan, 0, 0)
		stmt += "(?,?,?,?,?,?,?,?,?,?,?),"
		vals = append(vals, id, short, d.OriginalUrl, dic[d.Domain], d.Alias, now, now, expirationDate, organizationId, userId, d.Cloaking, false)
	}
	//trim the last ,
	stmt = stmt[0 : len(stmt)-1]
	return stmt, vals
}
