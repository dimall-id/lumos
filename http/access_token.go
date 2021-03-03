package http

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type AccessToken struct {
	Jti string `json:"jti" gorm:"jti;size:50"`
	UserId string `json:"user_id" gorm:"user_id;size:50"`
	RefId string `json:"ref_id" gorm:"ref_id;size:50"`
	Email string `json:"email" gorm:"email;size:255"`
	PhoneNo string `json:"phone_no" gorm:"phone_no;size:255"`
	UserName string `json:"user_name" gorm:"user_name:size:255"`
	UserType string `json:"user_type" gorm:"user_type;size:255"`
	Roles []string `json:"roles" gorm:"roles;type:varchar[]"`
	Iat int64 `json:"iat" gorm:"iat"`
	Exp int64 `json:"exp" gorm:"exp"`
}

func (u *AccessToken) BeforeCreate(tx *gorm.DB) (err error) {
	u.Jti = uuid.New().String()
	createdAt := time.Now().Unix()
	u.Iat = createdAt
	return
}

func (a *AccessToken) FillAccessToken (data map[string]interface{}) {
	if val, oke := data["exp"];oke {
		a.Exp = val.(int64)
	}
	if val, oke := data["iat"];oke {
		a.Iat = val.(int64)
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