package builder

import (
	"fmt"
	"github.com/dimall-id/lumos/v2/misc"
	"gorm.io/gorm"
	"regexp"
	"strconv"
	"strings"
	"time"
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
		split := strings.Split(cond["val_one"], "-")
		year, _ := strconv.Atoi(split[2])
		month, _ := strconv.Atoi(split[1])
		date, _ := strconv.Atoi(split[0])
		valOne := time.Date(year, time.Month(month), date, 0, 0, 0, 0, time.UTC)
		query := fmt.Sprintf("%s%s%d", field, GetOperator(cond["op_one"]), valOne.Unix())
		tx := tx.Where(query)
		return tx
	} else {
		split := strings.Split(cond["val_one"], "-")
		year, _ := strconv.Atoi(split[2])
		month, _ := strconv.Atoi(split[1])
		date, _ := strconv.Atoi(split[0])
		valOne := time.Date(year, time.Month(month), date, 0, 0, 0, 0, time.UTC)
		queryOne := fmt.Sprintf("%s%s%d", field, GetOperator(cond["op_one"]), valOne.Unix())

		split = strings.Split(cond["val_two"], "-")
		year, _ = strconv.Atoi(split[2])
		month, _ = strconv.Atoi(split[1])
		date, _ = strconv.Atoi(split[0])
		valTwo := time.Date(year, time.Month(month), date, 0, 0, 0, 0, time.UTC)
		queryTwo := fmt.Sprintf("%s%s%d", field, GetOperator(cond["op_two"]), valTwo.Unix())
		tx := tx.Where(queryOne).Where(queryTwo)
		return tx
	}
}
