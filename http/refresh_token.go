package http

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type RefreshToken struct {
	Rti string `json:"rti" gorm:"rti;size:50"`
	Jti string 	`json:"-" gorm:"jti;size:50"`
	UserId string `json:"user_id" gorm:"user_id;size:50"`
	IsClaimed bool `json:"-" gorm:"is_claimed"`
	Iat float64 `json:"iat" gorm:"iat"`
	Exp float64 `json:"exp" gorm:"exp"`
}

func (u *RefreshToken) BeforeCreate(tx *gorm.DB) (err error) {
	u.Rti = uuid.New().String()
	createdAt := time.Now().Unix()
	u.Iat = float64(createdAt)
	return
}

func (r *RefreshToken) FillRefreshToken (data map[string]interface{}) {
	if val, oke := data["iat"];oke {
		r.Iat = val.(float64)
	}
	if val, oke := data["exp"];oke {
		r.Exp = val.(float64)
	}
	if val, oke := data["jti"];oke {
		r.Jti = val.(string)
	}
	if val, oke := data["rti"];oke {
		r.Rti = val.(string)
	}
	if val,oke := data["user_id"];oke {
		r.UserId = val.(string)
	}
}