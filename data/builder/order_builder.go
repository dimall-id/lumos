package builder

import (
	"fmt"
	"github.com/dimall-id/lumos/misc"
	"gorm.io/gorm"
	"regexp"
	"strings"
)

const (
	OrderPattern = "\\[(?P<type>order);(?P<condition>(?:(?:[a-zA-Z\\_\\.]+):(?:asc|desc))(?:,?(?:[a-zA-Z\\_\\.]+):(?:asc|desc))*)\\]"
)

type OrderBuilder struct {}

func (ob *OrderBuilder) IsValid (value string) bool {
	r := regexp.MustCompile(OrderPattern)
	return r.MatchString(value)
}

func (ob *OrderBuilder) ApplyQuery (db *gorm.DB, field string, condition string) *gorm.DB {
	cond := misc.BuildToMap(OrderPattern, condition)
	if cond == nil {
		return db
	}
	orders := strings.Split(cond["condition"],",")
	tx := db
	for _, order := range orders {
		o := strings.Split(order, ":")
		query := fmt.Sprintf("%s %s", o[0], o[1])
		tx = tx.Order(query)
	}
	return tx
}