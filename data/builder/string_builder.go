package builder

import (
	"fmt"
	"github.com/dimall-id/lumos/v2/misc"
	log "github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"regexp"
	"strings"
)

const (
	StringPattern = "\\[(?P<type>(?:eq|neq|like|ilike));(?P<condition>[a-zA-Z0-9\\s\\-]+)\\]"
)

type StringBuilder struct {}


func (lb *StringBuilder) IsValid (value string) bool {
	r := regexp.MustCompile(StringPattern)
	return r.MatchString(value)
}

func (lb *StringBuilder) getCondition(operator string, condition string) string {
	if operator == "like" || operator == "ilike" {
		cond := strings.Replace(condition," ", "%", -1)
		c := fmt.Sprintf("%%%s%%%", cond)
		log.Error(c)
		return c
	}
	return condition
}

func (lb *StringBuilder) ApplyQuery (db *gorm.DB, field string, condition string) *gorm.DB {
	cond := misc.BuildToMap(StringPattern, condition)
	if cond == nil {
		return db
	}
	query := field + GetOperator(cond["type"]) + "'" + lb.getCondition(cond["type"], cond["condition"]) + "'"
	tx := db
	tx = tx.Where(query)
	return tx
}

