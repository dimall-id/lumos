package data

import (
	"fmt"
	"github.com/dimall-id/lumos/data/builder"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"testing"
)

func TestQuery_BuildData(t *testing.T) {
	testSets := []struct {
		queries map[string]string
		expect string
	}{
		{
			queries : map[string]string{
				"name" : "[like;%edited%]",
			},
			expect: "SELECT * FROM \"products\" WHERE (price > 3000 AND price <= 15000) AND name LIKE '%edit%' AND id != '31a01baf-b650-45a1-bb2e-211c6533518c'",
		},
		{
			queries : map[string]string{
				"select" : "[select;name,description]",
				"price" : "[numeric;gt:3000,lte:15000]",
			},
			expect: "SELECT \"name\",\"description\" FROM \"products\" WHERE price > 3000 AND price <= 15000",
		},
		{
			queries : map[string]string{
				"select" : "[select;name,description]",
				"id" : "[in;5138fffb-6d0b-4384-93b6-489aa890950b,31a01baf-b650-45a1-bb2e-211c6533518c]",
			},
			expect: "SELECT \"name\",\"description\" FROM \"products\" WHERE id IN ($1,$2)",
		},
	}

	var db, _ = gorm.Open(postgres.Open("host=localhost user=postgres password=postgres dbname=products port=5432 TimeZone=UTC"), &gorm.Config{
		//DryRun: true,
	})
	for _, test := range testSets {
		var data []builder.Product
		tx := db
		tx.Model(builder.Product{})
		Q := New(tx)
		res := Q.BuildResponse(test.queries, &data)
		fmt.Println(res.Data)
		t.Error(res.Data)
		//if strings.TrimSpace(tx.Statement.SQL.String()) != test.expect {
		//	t.Errorf("[%d] Fail to Test, Didn't generate SQL expected", i)
		//}
	}
}