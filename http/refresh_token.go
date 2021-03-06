package http

import (
	"github.com/dimall-id/jwt-go"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type RefreshToken struct {
	Rti string `json:"rti" gorm:"rti;primaryKey"`
	Jti string 	`json:"-" gorm:"jti;primaryKey"`
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

func (r *RefreshToken) FromJwtBase64 (base64Token string) error {
	token, err := jwt.ParseUnverified(base64Token, jwt.MapClaims{})
	if err != nil {
		return err
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	r.FillRefreshToken(claims)
	return nil
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