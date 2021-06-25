package builder

func GetOperator (op string) string {
	switch op {
	case "gt" :
		return " > "
	case "lt" :
		return " < "
	case "gte" :
		return " >= "
	case "lte" :
		return " <= "
	case "eq" :
		return " = "
	case "neq" :
		return " != "
	case "like" :
		return " LIKE "
	case "ilike" :
		return " ILIKE "
	}
	return " = "
}
