package misc

import "testing"

func TestBuildToMap(t *testing.T) {
	results := BuildToMap(`\[(?P<type>select);(?P<condition>[a-zA-Z,]+)\]`, `[select;name,description]`)
	if len(results) != 3 {
		t.Error("Fail to test, len of map should 3")
	}
	if results["condition"] != "name,description" {
		t.Error("Fail to test, value doesn't fit what is should be")
	}
}