package builder

import "testing"

func TestGetOperator(t *testing.T) {
	testSets := []struct{
		value string
		expect string
	}{
		{
			value: "gt",
			expect: " > ",
		},
		{
			value: "gte",
			expect: " >= ",
		},
		{
			value: "lt",
			expect: " < ",
		},
		{
			value: "lte",
			expect: " <= ",
		},
		{
			value: "eq",
			expect: " = ",
		},
		{
			value: "neq",
			expect: " != ",
		},
		{
			value: "unknow",
			expect: " = ",
		},
	}

	for i, test := range testSets {
		if GetOperator(test.value) != test.expect {
			t.Errorf("[%d] Fail to test, Get Operator method doesn't return result as expected [%s]", i, test.expect)
		}
	}
}
