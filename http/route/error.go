package route

import "fmt"

type InvalidRouteError struct {
	route Route
}

func (ir *InvalidRouteError) Error() string {
	var funcStatus string
	if ir.route.Func == nil {
		funcStatus = "NOT PARSED"
	} else {
		funcStatus = "PARSED"
	}
	return fmt.Sprintf("Route given is invalid. Http Method : %s, Url : %s, Func : %s", ir.route.HttpMethod, ir.route.Url, funcStatus)
}

type ExistingRouteError struct {
	route Route
}

func (er *ExistingRouteError) Error() string {
	return fmt.Sprintf("Existing Route with Http Method : %s and Url : %s existed", er.route.HttpMethod, er.route.Url)
}
