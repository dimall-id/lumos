package builder

import (
	"gorm.io/gorm"
	"regexp"
	"strings"
)
import "github.com/dimall-id/lumos/v2/misc"

const (
	SelectPattern = "\\[(?P<type>select):(?P<condition>[a-zA-Z,]+)\\]"
)

type SelectBuilder struct{}

func (sb *SelectBuilder) IsValid(value string) bool {
	r := regexp.MustCompile(SelectPattern)
	return r.MatchString(value)
}

func (sb *SelectBuilder) ApplyQuery(db *gorm.DB, field string, condition string) *gorm.DB {
	cond := misc.BuildToMap(SelectPattern, condition)
	if cond == nil {
		return db
	}
	fields := strings.Split(cond["condition"], ",")
	tx := db
	tx = tx.Select(fields)
	return tx
}
