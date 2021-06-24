package data

import (
	"fmt"
	"github.com/dimall-id/lumos/v2/data/builder"
	"github.com/dimall-id/lumos/v2/misc"
	"gorm.io/gorm"
	"net/url"
	"regexp"
	"strconv"
)

type ExistingQueryBuilderError struct {
	key string
}
func (e *ExistingQueryBuilderError) Error() string {
	return fmt.Sprintf("Existing Query Builder with key '%s' is found", e.key)
}

type NoExistingQueryBuilderError struct {
	key string
}
func (e *NoExistingQueryBuilderError) Error() string {
	return fmt.Sprintf("No Existing Query Builder with key '%s' is found", e.key)
}

func ExtractQuery (queries string) map[string]string {
	qString, _ := url.QueryUnescape(queries)
	r := regexp.MustCompile(`(?:(?P<key>[\w\d\_]+)=(?P<value>[\w\d\:\[\]\,\;\_\%.\-]+))+`)
	exps := r.FindAllStringSubmatch(queries, -1)
	var results = make(map[string]string)
	var keys = make(map[string]int)
	for i, key := range r.SubexpNames() {
		if key != "" {
			keys[key] = i
		}
	}
	for _, exp := range exps {
		t, err := url.QueryUnescape(exp[keys["value"]])
		if  err == nil {
		results[exp[keys["key"]]] = t
		}
	}
	return results
}

type QueryBuilder interface {
	IsValid (value string) bool
	ApplyQuery (db *gorm.DB, field string, condition string) *gorm.DB
}

type Query struct {
	db *gorm.DB
	builders map[string]QueryBuilder
}

func New (db *gorm.DB) Query  {
	datas := make(map[string]QueryBuilder)
	datas["date"] = &builder.DateBuilder{}
	datas["list"] = &builder.ListBuilder{}
	datas["numeric"] = &builder.NumericBuilder{}
	datas["order"] = &builder.OrderBuilder{}
	datas["select"] = &builder.SelectBuilder{}
	datas["string"] = &builder.StringBuilder{}
	datas["with"] = &builder.WithBuilder{}
	tx := db
	return Query{
		db: tx,
		builders: datas,
	}
}

func (q *Query) AddBuilder (key string, builder QueryBuilder) error {
	if _, oke := q.builders[key]; oke {
		return &ExistingQueryBuilderError{key: key}
	}
	q.builders[key] = builder
	return nil
}

func (q *Query) RemoveBuilder (key string) error {
	if _, oke := q.builders[key]; !oke {
		return &NoExistingQueryBuilderError{key: key}
	}
	delete(q.builders, key)
	return nil
}

func (q *Query) GetBuilder (value string) QueryBuilder {
	for _, b := range q.builders {
		if b.IsValid(value) {
			return b
		}
	}
	return nil
}

func (q *Query) BuildResponse (queries map[string]string, result interface{}) builder.HttpResponse {
	for field, condition := range queries {
		b := q.GetBuilder(condition)
		if b != nil {
			q.db = b.ApplyQuery(q.db, field, condition)
		}
	}

	if val, ok := queries["paging"]; ok {
		if builder.IsPagingValid(val) {
			paging := misc.BuildToMap(builder.PagingPattern, val)
			page,_ := strconv.ParseInt(paging["page"], 10, 32)
			perPage,_ := strconv.ParseInt(paging["per_page"], 10, 32)
			return *builder.Paging(&builder.Param{
				DB: q.db,
				Page: int(page),
				PerPage: int(perPage),
				ShowSQL: true,
			}, result)
		}
	}

	var count int64
	q.db.Model(result).Count(&count)
	q.db.Find(result)
	return builder.HttpResponse{
		Records: result,
		TotalRecord: int(count),
	}
}
