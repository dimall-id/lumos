package builder

import (
	"github.com/dimall-id/lumos/misc"
	"gorm.io/gorm"
	"regexp"
)

const (
	StringPattern = "\\[(?P<type>(?:eq|neq|like)):(?P<condition>[a-zA-Z0-9\\s\\%\\-]+)\\]"
)

type StringBuilder struct{}

func (lb *StringBuilder) IsValid(value string) bool {
	r := regexp.MustCompile(StringPattern)
	return r.MatchString(value)
}

func (lb *StringBuilder) ApplyQuery(db *gorm.DB, field string, condition string) *gorm.DB {
	cond := misc.BuildToMap(StringPattern, condition)
	if cond == nil {
		return db
	}
	query := field + GetOperator(cond["type"]) + "'" + cond["condition"] + "'"
	tx := db
	tx = tx.Where(query)
	return tx
}
