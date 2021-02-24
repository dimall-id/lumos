package http

type AccessToken struct {
	Iss string `json:"iss"`
	Exp string `json:"exp"`
	Jti string `json:"jti"`
	UserId string `json:"user_id"`
	UserName string `json:"user_name"`
	UserType string `json:"user_type"`
	Roles []string `json:"roles"`
}

func (a *AccessToken) FillAccessToken (data map[string]interface{}) {
	if val, oke := data["iss"];oke {
		a.Iss = val.(string)
	}
	if val, oke := data["exp"];oke {
		a.Exp = val.(string)
	}
	if val, oke := data["jti"];oke {
		a.Jti = val.(string)
	}
	if val,oke := data["user_id"];oke {
		a.UserId = val.(string)
	}
	if val,oke := data["user_name"];oke {
		a.UserName = val.(string)
	}
	if val,oke := data["user_type"];oke {
		a.UserType = val.(string)
	}
	if val,oke := data["roles"];oke {
		a.Roles = val.([]string)
	}
}