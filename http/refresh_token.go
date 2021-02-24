package http

type RefreshToken struct {
	Iss string `json:"iss"`
	Exp string `json:"exp"`
	Jti string `json:"jti"`
	UserId string `json:"user_id"`
}

func (r *RefreshToken) FillRefreshToken (data map[string]interface{}) {
	if val, oke := data["iss"];oke {
		r.Iss = val.(string)
	}
	if val, oke := data["exp"];oke {
		r.Exp = val.(string)
	}
	if val, oke := data["jti"];oke {
		r.Jti = val.(string)
	}
	if val,oke := data["user_id"];oke {
		r.UserId = val.(string)
	}
}