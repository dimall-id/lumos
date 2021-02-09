package builder

import (
	"github.com/dimall-id/lumos/misc"
	"gorm.io/gorm"
	"regexp"
	"strings"
)

const (
	WithPattern = "\\[(?P<type>with);(?P<condition>[a-zA-Z,]+)\\]"
)

type WithBuilder struct {}

func (wb *WithBuilder) IsValid (value string) bool {
	r := regexp.MustCompile(WithPattern)
	return r.MatchString(value)
}

func (wb *WithBuilder) ApplyQuery (db *gorm.DB, field string, condition string) *gorm.DB {
	cond := misc.BuildToMap(WithPattern, condition)
	if cond == nil {
		return db
	}
	relations := strings.Split(cond["condition"], ",")
	tx := db
	for _, relation := range relations {
		tx = db.Preload(relation)
	}
	return tx
}