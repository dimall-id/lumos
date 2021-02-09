package builder

import (
	"gorm.io/gorm"
	"math"
	"regexp"
)

type HttpResponse struct {
	CurrentPage int `json:"current_page,omitempty"`
	LastPageNumber int `json:"last_page_number,omitempty"`
	Total int `json:"total,omitempty"`
	PerPage int `json:"per_page,omitempty"`
	Data interface{} `json:"data"`
	PrevPage string `json:"prev_page,omitempty"`
	NextPage string `json:"nex_page,omitempty"`
	FirstPage string `json:"first_page,omitempty"`
	LastPage string `json:"last_page,omitempty"`
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

	results.Total = int(count)
	results.Data = result
	results.CurrentPage = p.Page
	results.LastPageNumber = int(math.Ceil(float64(count) / float64(p.PerPage)))

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