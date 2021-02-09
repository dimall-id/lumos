package data

import (
	"github.com/dimall-id/lumos/data/builder"
	"gorm.io/gorm"
)

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

func (q *Query) AddBuilder (key string, builder QueryBuilder) {
	q.builders[key] = builder
}

func (q *Query) RemoveBuilder (key string) {
	delete(q.builders, key)
}

func (q *Query) GetBuilder (value string) QueryBuilder {
	for _, b := range q.builders {
		if b.IsValid(value) {
			return b
		}
	}
	return nil
}

func (q *Query) BuildResponse (queries map[string]string, result interface{}, path string) builder.HttpResponse {
	for field, condition := range queries {
		b := q.GetBuilder(condition)
		if b != nil {
			q.db = b.ApplyQuery(q.db, field, condition)
		}
	}

	q.db.Find(result)
	res := builder.HttpResponse{
		Data: result,
	}
	return res
}

