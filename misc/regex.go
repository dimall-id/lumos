package misc

import (
	"regexp"
)

func BuildToMap (pattern string, value string) map[string]string {
	r := regexp.MustCompile(pattern)
	values := r.FindStringSubmatch(value)
	if len(values) >= len(r.SubexpNames()) {
		results := make(map[string]string)
		for i, key := range r.SubexpNames() {
			if key == "" {
				results["raw"] = values[i]
			} else {
				results[key] = values[i]
			}
		}
		return results
	}
	return nil
}