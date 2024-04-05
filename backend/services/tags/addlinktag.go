package services

import (
	"errors"
	"net/http"
	"strings"
	"time"

	sequentialguid "github.com/adewoleadenigbagbe/sequential-guid"
	"github.com/adewoleadenigbagbe/url-shortner-service/models"
	"github.com/labstack/echo/v4"
	"github.com/samber/lo"
)

type shortTagInfo struct {
	Id        string
	TagId     string
	ShortId   string
	CreatedOn time.Time
}

func (service TagService) AddShortLinkTag(tagContext echo.Context) error {
	var (
		err error
	)

	request := new(models.CreateShortLinkTagRequest)
	err = tagContext.Bind(request)
	if err != nil {
		return tagContext.JSON(http.StatusBadRequest, []string{err.Error()})
	}

	modelErrs := validateShortLinkTagRequest(*request)
	if len(modelErrs) > 0 {
		valErrors := lo.Map(modelErrs, func(er error, index int) string {
			return er.Error()
		})
		return tagContext.JSON(http.StatusBadRequest, valErrors)
	}

	request.Tags = lo.Map(request.Tags, func(tag string, _ int) string {
		return strings.ToLower(tag)
	})
	var tagstoAdd []tagInfo
	var shortTagstoAdd []shortTagInfo

	tagQuery, tagArgs := formatTagQuery(request.Tags)
	rows, err := service.Db.Query(tagQuery, tagArgs...)
	if err != nil {
		return tagContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	defer rows.Close()

	var existingTagInfos []tagInfo
	for rows.Next() {
		var existingTagInfo tagInfo
		err = rows.Scan(&existingTagInfo.Id, &existingTagInfo.Name)
		if err != nil {
			return tagContext.JSON(http.StatusInternalServerError, []string{err.Error()})
		}
		existingTagInfos = append(existingTagInfos, existingTagInfo)
	}

	if len(existingTagInfos) > 0 {
		if len(request.Tags) > len(existingTagInfos) {
			lo.ForEach(request.Tags, func(tagName string, _ int) {
				doesExist := lo.ContainsBy(existingTagInfos, func(existingTag tagInfo) bool {
					return existingTag.Name == tagName
				})

				if !doesExist {
					newtagInfo := tagInfo{
						Id:   sequentialguid.NewSequentialGuid().String(),
						Name: tagName,
					}
					tagstoAdd = append(tagstoAdd, newtagInfo)

					newshortTagInfo := shortTagInfo{
						Id:        sequentialguid.NewSequentialGuid().String(),
						TagId:     newtagInfo.Id,
						ShortId:   request.ShortId,
						CreatedOn: time.Now(),
					}
					shortTagstoAdd = append(shortTagstoAdd, newshortTagInfo)
				}
			})
		}

		linkTagQuery, linkQueryArgs := formatShortlinkTagQuery(existingTagInfos, request.ShortId)
		queryRows, err := service.Db.Query(linkTagQuery, linkQueryArgs...)
		if err != nil {
			return tagContext.JSON(http.StatusInternalServerError, []string{err.Error()})
		}

		defer queryRows.Close()

		var shortTagInfos []shortTagInfo
		for queryRows.Next() {
			var shortTagInfo shortTagInfo
			err = queryRows.Scan(&shortTagInfo.Id, &shortTagInfo.ShortId, &shortTagInfo.TagId, &shortTagInfo.CreatedOn)
			if err != nil {
				return tagContext.JSON(http.StatusInternalServerError, []string{err.Error()})
			}
			shortTagInfos = append(shortTagInfos, shortTagInfo)
		}

		if len(shortTagInfos) > 0 {
			if len(existingTagInfos) > len(shortTagInfos) {
				lo.ForEach(existingTagInfos, func(existingTag tagInfo, _ int) {
					doesNotExist := lo.ContainsBy(shortTagInfos, func(existingShortTagInfo shortTagInfo) bool {
						return existingShortTagInfo.TagId == existingTag.Id
					})

					if !doesNotExist {
						shortTagInfo := shortTagInfo{
							Id:        sequentialguid.NewSequentialGuid().String(),
							TagId:     existingTag.Id,
							ShortId:   request.ShortId,
							CreatedOn: time.Now(),
						}
						shortTagstoAdd = append(shortTagstoAdd, shortTagInfo)
					}
				})
			}
		} else {
			shortTags := lo.Map(existingTagInfos, func(existingTag tagInfo, _ int) shortTagInfo {
				return shortTagInfo{
					Id:        sequentialguid.NewSequentialGuid().String(),
					TagId:     existingTag.Id,
					ShortId:   request.ShortId,
					CreatedOn: time.Now(),
				}
			})

			shortTagstoAdd = append(shortTagstoAdd, shortTags...)
		}
	} else {
		tagstoAdd = lo.Map(request.Tags, func(tag string, _ int) tagInfo {
			return tagInfo{
				Id:   sequentialguid.NewSequentialGuid().String(),
				Name: tag,
			}
		})

		shortTagstoAdd = lo.Map(tagstoAdd, func(tagInfo tagInfo, _ int) shortTagInfo {
			return shortTagInfo{
				Id:        sequentialguid.NewSequentialGuid().String(),
				TagId:     tagInfo.Id,
				ShortId:   request.ShortId,
				CreatedOn: time.Now(),
			}
		})
	}

	tx, err := service.Db.Begin()
	if err != nil {
		return tagContext.JSON(http.StatusInternalServerError, []string{err.Error()})
	}

	if len(tagstoAdd) > 0 {
		stmt, args := insertTagStmt(tagstoAdd)
		_, err = tx.Exec(stmt, args...)
		if err != nil {
			tx.Rollback()
			return tagContext.JSON(http.StatusInternalServerError, []string{err.Error()})
		}
	}

	if len(shortTagstoAdd) > 0 {
		shortLinkStmt, linkArgs := insertShortlinkTagStmt(shortTagstoAdd)
		_, err = tx.Exec(shortLinkStmt, linkArgs...)
		if err != nil {
			tx.Rollback()
			return tagContext.JSON(http.StatusInternalServerError, []string{err.Error()})
		}
	}
	tx.Commit()
	return tagContext.JSON(http.StatusCreated, nil)
}

func validateShortLinkTagRequest(request models.CreateShortLinkTagRequest) []error {
	var validationErrors []error

	if request.ShortId == "" {
		validationErrors = append(validationErrors, errors.New("shortId is required"))
	}

	if len(request.Tags) == 0 {
		validationErrors = append(validationErrors, errors.New("tags is empty. supply at least one"))
	}
	return validationErrors
}

func formatShortlinkTagQuery(tagInfos []tagInfo, shortId string) (string, []interface{}) {
	str := "SELECT * FROM shortlinktags WHERE TagId IN (?" + strings.Repeat(",?", len(tagInfos)-1) + ") AND ShortId =?"
	vals := []interface{}{}
	for _, tagInfo := range tagInfos {
		vals = append(vals, tagInfo.Id)
	}
	vals = append(vals, shortId)
	return str, vals
}

func insertShortlinkTagStmt(shortTagInfos []shortTagInfo) (string, []interface{}) {
	stmt := "INSERT INTO shortlinktags VALUES "
	vals := []interface{}{}

	for _, shortTagInfo := range shortTagInfos {
		stmt += "(?,?,?,?),"
		vals = append(vals, shortTagInfo.Id, shortTagInfo.ShortId, shortTagInfo.TagId, shortTagInfo.CreatedOn)
	}
	//trim the last ,
	stmt = stmt[0 : len(stmt)-1]
	return stmt, vals
}
