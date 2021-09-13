package builder

import (
	"fmt"
	"github.com/dimall-id/lumos/v2/misc"
	"gorm.io/gorm"
	"regexp"
	"time"
)

const (
	DatePattern = "\\[(?P<type>date);(?P<condition>(?:(?P<op_one>gt|gte|eq|neq):(?P<val_one>\\d{2}-\\d{2}-\\d{4}))?,?(?:(?P<op_two>lt|lte):(?P<val_two>\\d{2}-\\d{2}-\\d{4}))?)\\]"
)

type DateBuilder struct {}

func (dd *DateBuilder) IsValid (value string) bool {
	r := regexp.MustCompile(DatePattern)
	return r.MatchString(value)
}

func (dd *DateBuilder) ApplyQuery (db *gorm.DB, field string, condition string) *gorm.DB {
	cond := misc.BuildToMap(DatePattern, condition)
	if cond == nil {return db}
	tx := db
	format := "12-30-1993"
	if cond["op_two"] == "" || cond["op_one"] == "eq" || cond["op_one"] == "neq" {
		valOne, _ := time.Parse(format, cond["val_one"])
		query := fmt.Sprintf("%s%s'%d'", field, GetOperator(cond["op_one"]), valOne.Unix())
		tx := tx.Where(query)
		return tx
	} else {
		valOne, _ := time.Parse(format, cond["val_one"])
		queryOne := fmt.Sprintf("%s%s'%d'", field, GetOperator(cond["op_one"]), valOne.Unix())
		valTwo, _ := time.Parse(format, cond["val_two"])
		queryTwo := fmt.Sprintf("%s%s'%d'", field, GetOperator(cond["op_one"]), valTwo.Unix())
		tx := tx.Where(queryOne).Where(queryTwo)
		return tx
	}
}