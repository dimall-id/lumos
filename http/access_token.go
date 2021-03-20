package http

import (
	"github.com/dimall-id/jwt-go"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
	"time"
)

type AccessToken struct {
	Jti string `json:"jti" gorm:"jti;primaryKey"`
	UserId string `json:"user_id" gorm:"user_id"`
	RefId string `json:"ref_id" gorm:"ref_id"`
	Email string `json:"email" gorm:"email"`
	PhoneNo string `json:"phone_no" gorm:"phone_no"`
	UserName string `json:"user_name" gorm:"user_name"`
	UserType string `json:"user_type" gorm:"user_type"`
	Roles pq.StringArray `json:"roles" gorm:"roles"`
	Iat float64 `json:"iat" gorm:"iat"`
	Exp float64 `json:"exp" gorm:"exp"`
}

func (a *AccessToken) BeforeCreate(tx *gorm.DB) (err error) {
	a.Jti = uuid.New().String()
	createdAt := time.Now().Unix()
	a.Iat = float64(createdAt)
	return
}

func (a *AccessToken) FromJwtBase64 (base64Token string) error {
	token, err := jwt.ParseUnverified(base64Token, jwt.MapClaims{})
	if err != nil {
		return err
	}
	claims, _ := token.Claims.(jwt.MapClaims)
	a.FillAccessToken(claims)
	return nil
}

func (a *AccessToken) FillAccessToken (data map[string]interface{}) {
	if val, oke := data["exp"];oke {
		a.Exp = val.(float64)
	}
	if val, oke := data["iat"];oke {
		a.Iat = val.(float64)
	}
	if val, oke := data["jti"];oke {
		a.Jti = val.(string)
	}
	if val,oke := data["user_id"];oke {
		a.UserId = val.(string)
	}
	if val,oke := data["ref_id"];oke {
		a.RefId = val.(string)
	}
	if val,oke := data["user_name"];oke {
		a.UserName = val.(string)
	}
	if val,oke := data["user_type"];oke {
		a.UserType = val.(string)
	}
	if val,oke := data["roles"];oke {
		roles := make([]string, len(val.([]interface{})))
		for i, role := range data["roles"].([]interface{}) {
			roles[i] = role.(string)
		}
		a.Roles = roles
	}
}