package builder

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"strings"
	"testing"
)

func TestSelectBuilder_IsValid(t *testing.T) {
	testSets := []struct {
		value string
		expect bool
	}{
		{
			"[select;name,description]",
			true,
		},
		{
			"[gte:2020-02-20]",
			false,
		},
	}

	for i,test := range testSets {
		sb := SelectBuilder{}
		isValid := sb.IsValid (test.value)
		if isValid != test.expect {
			t.Errorf("[%d] Fail to test, result doesn't meet expectation", i)
		}
	}
}

func TestSelectBuilder_AddQuery(t *testing.T) {
	testSets := []struct {
		field string
		condition string
		expect string
	}{
		{
			"select",
			"[select;name,description]",
			`SELECT "name","description" FROM "products"`,
		},
	}
	var db, _ = gorm.Open(postgres.Open("host=localhost user=postgres password=postgres dbname=products port=5432 TimeZone=UTC"), &gorm.Config{
		DryRun: true,
	})
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	for _, test := range testSets {
		sb := SelectBuilder{}
		db = sb.ApplyQuery(db, test.field, test.condition)
		var datas []Product
		stmt := db.Find(&datas).Statement;
		if strings.TrimSpace(stmt.SQL.String()) != test.expect {
			t.Error("Fail to test, doesn't generate expected SQL query")
		}
	}
}
