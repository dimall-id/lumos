package builder

import (
	"github.com/dimall-id/lumos/misc"
	"gorm.io/gorm"
	"regexp"
)

const (
	DatePattern = "\\[(?:(?P<op_one>gt|gte|eq|neq):(?P<val_one>\\d{2}-\\d{2}-\\d{4}))?,?(?:(?P<op_two>lt|lte):(?P<val_two>\\d{2}-\\d{2}-\\d{4}))?\\]"
)

type DateBuilder struct{}

func (dd *DateBuilder) IsValid(value string) bool {
	r := regexp.MustCompile(DatePattern)
	return r.MatchString(value)
}

func (dd *DateBuilder) ApplyQuery(db *gorm.DB, field string, condition string) *gorm.DB {
	cond := misc.BuildToMap(DatePattern, condition)
	if cond == nil {
		return db
	}
	tx := db
	if cond["op_two"] == "" || cond["op_one"] == "eq" || cond["op_one"] == "neq" {
		query := field + GetOperator(cond["op_one"]) + "'" + cond["val_one"] + "'"
		tx := tx.Where(query)
		return tx
	} else {
		queryOne := field + GetOperator(cond["op_one"]) + "'" + cond["val_one"] + "'"
		queryTwo := field + GetOperator(cond["op_two"]) + "'" + cond["val_two"] + "'"
		tx := tx.Where(queryOne).Where(queryTwo)
		return tx
	}
}
