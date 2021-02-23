package builder

import (
	"gorm.io/gorm"
	"math"
	"regexp"
)

type HttpResponse struct {
	TotalRecord int `json:"total_record,omitempty" msgpack:"total_record,omitempty"`
	TotalPage int `json:"total_page,omitempty" msgpack:"total_page,omitempty"`
	Page int `json:"page,omitempty" msgpack:"page,omitempty"`
	PerPage int `json:"per_page,omitempty" msgpack:"per_page,omitempty"`
	PrevPage int `json:"prev_page,omitempty" msgpack:"prev_page,omitempty"`
	NextPage int `json:"next_page,omitempty" msgpack:"next_page,omitempty"`
	FirstPage int `json:"first_page,omitempty" msgpack:"first_page,omitempty"`
	LastPage int `json:"last_page,omitempty" msgpack:"last_page,omitempty"`
	Records interface{} `json:"records" msgpack:"records,as_array"`
}

type Param struct {
	DB *gorm.DB
	Page int
	PerPage int
	Path string
	ShowSQL bool
}

const (
	PagingPattern = "\\[page:(?P<page>[\\d]+),per_page:(?P<per_page>[\\d]+)\\]"
)

func Paging (p *Param, result interface{}) *HttpResponse {
	db := p.DB

	if p.ShowSQL {
		db = db.Debug()
	}
	if p.Page < 1 {
		p.Page = 1
	}
	if p.PerPage == 0 {
		p.PerPage = 10
	}

	done := make(chan bool, 1)
	var results HttpResponse
	var count int64
	var offset int

	go countRecords(db, result, done, &count)

	if p.Page == 1 {
		offset = 0
	} else {
		offset = (p.Page - 1) * p.PerPage
	}

	db.Limit(p.PerPage).Offset(offset).Find(result)
	<-done

	results.TotalRecord = int(count)
	results.Records = result
	results.Page = p.Page
	results.TotalPage = int(math.Ceil(float64(count) / float64(p.PerPage)))
	results.LastPage = results.TotalPage
	results.FirstPage = 1
	if p.Page > 1 { results.PrevPage = p.Page - 1 }
	if p.Page < results.LastPage { results.NextPage = p.Page + 1 }

	return &results
}

func countRecords(db *gorm.DB, anyType interface{}, done chan bool, count *int64) {
	db.Model(anyType).Count(count)
	done <- true
}

func IsPagingValid(value string) bool {
	r := regexp.MustCompile(PagingPattern)
	return r.MatchString(value)
}