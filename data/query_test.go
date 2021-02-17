package data

import (
	"gorm.io/gorm"
	"gorm.io/gorm/utils/tests"
	"testing"
)

func TestNew(t *testing.T) {
	var db, _ = gorm.Open(tests.DummyDialector{}, &gorm.Config{
		DryRun: true,
	})
	q := New(db)
	if _, oke := q.builders["date"]; !oke {
		t.Error("Fail to test, Query builder doesn't generate default Query Builder")
	}
}

func TestQuery_Query(t *testing.T) {

}