package builder

import (
	"fmt"
	"github.com/dimall-id/lumos/misc"
	"gorm.io/gorm"
	"regexp"
	"strings"
)

const (
	ListPattern = "\\[(?P<type>(?:in|nin));(?P<condition>[a-zA-Z0-9\\s\\%\\-\\,]+)\\]"
)

type ListBuilder struct {}

func (lb *ListBuilder) IsValid (value string) bool {
	r := regexp.MustCompile(ListPattern)
	return r.MatchString(value)
}

func (lb *ListBuilder) ApplyQuery (db *gorm.DB, field string, condition string) *gorm.DB {
	cond := misc.BuildToMap(ListPattern, condition)
	if cond == nil {
		return db
	}
	fields := strings.Split(cond["condition"], ",")
	var query string
	if strings.ToUpper(cond["type"]) == "IN" {
		query = fmt.Sprintf("%s IN ?", field)
	} else {
		query = fmt.Sprintf("%s NOT IN ?", field)
	}
	tx := db
	tx = tx.Where(query, fields)
	return tx
}
