package builder

import (
	"github.com/dimall-id/lumos/misc"
	"gorm.io/gorm"
	"regexp"
)

const (
	NumericPattern = "\\[(?:(?P<op_one>gt|gte|eq|neq):(?P<val_one>[\\d]+))?,?(?:(?P<op_two>lt|lte):(?P<val_two>[\\d]+))?\\]"
)

type NumericBuilder struct{}

func (dd *NumericBuilder) IsValid(value string) bool {
	r := regexp.MustCompile(NumericPattern)
	return r.MatchString(value)
}

func (dd *NumericBuilder) ApplyQuery(db *gorm.DB, field string, condition string) *gorm.DB {
	cond := misc.BuildToMap(NumericPattern, condition)
	if cond == nil {
		return db
	}
	tx := db
	if cond["op_two"] == "" || cond["op_one"] == "eq" || cond["op_one"] == "neq" {
		query := field + GetOperator(cond["op_one"]) + cond["val_one"]
		tx = tx.Where(query)
	} else {
		queryOne := field + GetOperator(cond["op_one"]) + cond["val_one"]
		queryTwo := field + GetOperator(cond["op_two"]) + cond["val_two"]
		tx = tx.Where(queryOne + " AND " + queryTwo)
	}
	return tx
}
