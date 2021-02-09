package builder

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
	"testing"
)

func TestListBuilder_IsValid(t *testing.T) {
	testSets := []struct {
		value string
		expect bool
	}{
		{
			"[in;name,description]",
			true,
		},
		{
			"[nin;name,description]",
			true,
		},
		{
			"[nin;name,Andy Wijaya]",
			true,
		},
	}

	for i,test := range testSets {
		sb := ListBuilder{}
		isValid := sb.IsValid (test.value)
		if isValid != test.expect {
			t.Errorf("[%d] Fail to test, result doesn't meet expectation", i)
		}
	}
}

func TestListBuilder_ApplyQuery(t *testing.T) {
	testSets := []struct {
		field string
		condition string
		expect string
	}{
		{
			"name",
			"[in;name,description]",
			`SELECT * FROM "products" WHERE name IN ($1,$2)`,
		},
		{
			"name",
			"[nin;name,description]",
			`SELECT * FROM "products" WHERE name NOT IN ($1,$2)`,
		},
	}
	for _, test := range testSets {
		var db, _ = gorm.Open(postgres.Open("host=localhost user=postgres password=postgres dbname=products port=5432 TimeZone=UTC"), &gorm.Config{
			DryRun: true,
		})
		sqlDB, _ := db.DB()
		sb := ListBuilder{}
		db = sb.ApplyQuery(db, test.field, test.condition)
		var datas []Product
		stmt := db.Find(&datas).Statement
		if strings.TrimSpace(stmt.SQL.String()) != test.expect {
			t.Error("Fail to test, doesn't generate expected SQL query")
		}
		defer sqlDB.Close()
	}
}
