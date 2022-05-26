package builder

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
	"testing"
)

func TestNumericBuilder_IsValid(t *testing.T) {
	testSets := []struct {
		value  string
		expect bool
	}{
		{
			"[gte:1000]",
			true,
		},
		{
			"[gte:1000,lte:2000]",
			true,
		},
		{
			"[lte:12000]",
			true,
		},
		{
			"[eq:12000]",
			true,
		},
		{
			"[neq:12000]",
			true,
		},
	}

	for i, test := range testSets {
		db := NumericBuilder{}
		isValid := db.IsValid(test.value)
		if isValid != test.expect {
			t.Errorf("[%d] Fail to test, result doesn't meet expectation", i)
		}
	}
}

func TestNumericBuilder_ApplyQuery(t *testing.T) {
	testSets := []struct {
		field     string
		condition string
		expect    string
	}{
		{
			"price",
			"[gte:1000,lt:2000]",
			"SELECT * FROM \"products\" WHERE price >= 1000 AND price < 2000",
		},
		{
			"price",
			"[gt:1500,lt:3000]",
			"SELECT * FROM \"products\" WHERE price > 1500 AND price < 3000",
		},
		{
			"price",
			"[gt:3000,lte:5000]",
			"SELECT * FROM \"products\" WHERE price > 3000 AND price <= 5000",
		},
		{
			"price",
			"[eq:5000]",
			"SELECT * FROM \"products\" WHERE price = 5000",
		},
		{
			"price",
			"[neq:5000]",
			"SELECT * FROM \"products\" WHERE price != 5000",
		},
	}
	for _, test := range testSets {
		var db, _ = gorm.Open(postgres.Open("host=localhost user=postgres password=postgres dbname=products port=5432 TimeZone=UTC"), &gorm.Config{
			DryRun: true,
		})
		sqlDB, _ := db.DB()
		dd := NumericBuilder{}
		field := test.field
		condition := test.condition
		db = dd.ApplyQuery(db, field, condition)
		var datas []Product
		stmt := db.Find(&datas).Statement
		if strings.TrimSpace(stmt.SQL.String()) != test.expect {
			t.Error("Fail to test, doesn't generate expected SQL query")
		}
		sqlDB.Close()
	}
}
