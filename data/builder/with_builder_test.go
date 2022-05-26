package builder

import "testing"

func TestWithBuilder_IsValid(t *testing.T) {
	testSets := []struct {
		value  string
		expect bool
	}{
		{
			"[with:productImages]",
			true,
		},
	}

	for i, test := range testSets {
		db := WithBuilder{}
		isValid := db.IsValid(test.value)
		if isValid != test.expect {
			t.Errorf("[%d] Fail to test, result doesn't meet expectation", i)
		}
	}
}
