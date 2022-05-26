package builder

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
	"testing"
)

func TestDateBuilder_IsValid(t *testing.T) {
	testSets := []struct {
		value  string
		expect bool
	}{
		{
			"[gte:23-02-2020]",
			true,
		},
		{
			"[gte:23-02-2020,lte:23-04-2020]",
			true,
		},
		{
			"[lt:23-04-2020,gt:23-02-2020]",
			false,
		},
		{
			"[lte:23-02-2020]",
			true,
		},
		{
			"[eq:23-02-2020]",
			true,
		},
		{
			"[neq:23-02-2020]",
			true,
		},
	}

	for i, test := range testSets {
		db := DateBuilder{}
		isValid := db.IsValid(test.value)
		if isValid != test.expect {
			t.Errorf("[%d] Fail to test, result doesn't meet expectation", i)
		}
	}
}

func TestDateBuilder_ApplyQuery(t *testing.T) {
	testSets := []struct {
		field     string
		condition string
		expect    string
	}{
		{
			"date",
			"[gte:23-10-2020,lt:24-11-2020]",
			"SELECT * FROM \"products\" WHERE date >= 1603411200 AND date < 1606176000",
		},
		{
			"date",
			"[gt:23-10-2020,lt:24-11-2020]",
			"SELECT * FROM \"products\" WHERE date > 1603411200 AND date < 1606176000",
		},
		{
			"date",
			"[gt:23-10-2020,lte:24-11-2020]",
			"SELECT * FROM \"products\" WHERE date > 1603411200 AND date <= 1606176000",
		},
		{
			"date",
			"[eq:23-10-2020]",
			"SELECT * FROM \"products\" WHERE date = 1603411200",
		},
		{
			"date",
			"[neq:23-10-2020]",
			"SELECT * FROM \"products\" WHERE date != 1603411200",
		},
	}
	for i, test := range testSets {
		var db, _ = gorm.Open(postgres.Open("host=localhost user=postgres password=postgres dbname=products port=5432 TimeZone=UTC"), &gorm.Config{
			DryRun: true,
		})
		sqlDB, _ := db.DB()
		dd := DateBuilder{}
		field := test.field
		condition := test.condition
		db = dd.ApplyQuery(db, field, condition)
		var datas []Product
		stmt := db.Find(&datas).Statement
		if strings.TrimSpace(stmt.SQL.String()) != test.expect {
			t.Errorf("[%d] Fail to test, doesn't generate expected SQL query", i)
		}
		sqlDB.Close()
	}
}
