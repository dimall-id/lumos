package builder

import (
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
	"testing"
)

func TestOrderBuilder_IsValid(t *testing.T) {
	testSets := []struct{
		value string
		expect bool
	}{
		{
			value: "[order;name:desc,id:asc]",
			expect: true,
		},
		{
			value: "[order;id:ASC]",
			expect: false,
		},
		{
			value: "[order;id:asc,name:desc,]",
			expect: false,
		},
		{
			value: "[order;id:desc]",
			expect: true,
		},
	}

	for i,test := range testSets {
		db := OrderBuilder{}
		isValid := db.IsValid (test.value)
		if isValid != test.expect {
			t.Errorf("[%d] Fail to test, result doesn't meet expectation", i)
		}
	}
}

func TestOrderBuilder_ApplyQuery(t *testing.T) {
	testSets := []struct{
		field string
		condition string
		expect string
	}{
		{
			field: "order",
			condition: "[order;id:desc,name:asc]",
			expect: `SELECT * FROM "products" ORDER BY id desc,name asc`,
		},
		{
			field: "order",
			condition: "[order;id:asc]",
			expect: `SELECT * FROM "products" ORDER BY id asc`,
		},
	}

	for i,test := range testSets {
		var db, _ = gorm.Open(postgres.Open("host=localhost user=postgres password=postgres dbname=products port=5432 TimeZone=UTC"), &gorm.Config{
			DryRun: true,
		})
		sqlDB, _ := db.DB()
		dd := OrderBuilder{}
		field := test.field
		condition := test.condition
		db = dd.ApplyQuery(db, field, condition)
		var datas []Product
		stmt := db.Find(&datas).Statement
		fmt.Println(stmt.SQL.String())
		if strings.TrimSpace(stmt.SQL.String()) != test.expect {
			t.Errorf("[%d] Fail to test, doesn't generate expected SQL query", i)
		}
		sqlDB.Close()
	}
}