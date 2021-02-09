package builder

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
	"testing"
)

func TestStringBuilder_IsValid(t *testing.T) {
	testSets := []struct {
		value string
		expect bool
	}{
		{
			"[like;%Andy Wijaya]",
			true,
		},
		{
			"[eq;Andy]",
			true,
		},
		{
			"[neq;Andy]",
			true,
		},
		{
			"[neq;Andy,Andy]",
			false,
		},
	}

	for i,test := range testSets {
		sb := StringBuilder{}
		isValid := sb.IsValid (test.value)
		if isValid != test.expect {
			t.Errorf("[%d] Fail to test, result doesn't meet expectation", i)
		}
	}
}

func TestStringBuilder_ApplyQuery(t *testing.T) {
	testSets := []struct {
		field string
		condition string
		expect string
	}{
		{
			"name",
			"[eq;Andy Wijaya]",
			`SELECT * FROM "products" WHERE name = 'Andy Wijaya'`,
		},
		{
			"name",
			"[neq;Andy Wijaya]",
			`SELECT * FROM "products" WHERE name != 'Andy Wijaya'`,
		},
		{
			"name",
			"[like;%Andy Wijaya%]",
			`SELECT * FROM "products" WHERE name LIKE '%Andy Wijaya%'`,
		},
		{
			"id",
			"[eq;31a01baf-b650-45a1-bb2e-211c6533518c]",
			`SELECT * FROM "products" WHERE id = '31a01baf-b650-45a1-bb2e-211c6533518c'`,
		},
	}
	for _, test := range testSets {
		var db, _ = gorm.Open(postgres.Open("host=localhost user=postgres password=postgres dbname=products port=5432 TimeZone=UTC"), &gorm.Config{
			DryRun: true,
		})
		sqlDB, _ := db.DB()
		sb := StringBuilder{}
		db = sb.ApplyQuery(db, test.field, test.condition)
		var datas []Product
		stmt := db.Find(&datas).Statement
		if strings.TrimSpace(stmt.SQL.String()) != test.expect {
			t.Error("Fail to test, doesn't generate expected SQL query")
		}
		sqlDB.Close()
	}
}