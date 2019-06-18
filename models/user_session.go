package models

import "time"

// LoginMethod 登录方式
type LoginMethod int32

const (
	// LoginMethodAccountPassword 用户密码登录
	LoginMethodAccountPassword LoginMethod = 1
	// LoginMethodPhoneVerifyCode 手机验证码登录
	LoginMethodPhoneVerifyCode LoginMethod = 2
	// LoginMethodPhonePassword 手机号密码登录
	LoginMethodPhonePassword LoginMethod = 3
)

// UserSession 客户登录记录
type UserSession struct {
	ID     uint64 `gorm:"primary_key" json:"id" form:"id" valid:"null"`
	UserID uint64 `gorm:"index" json:"user_id" valid:"null"`
	// User   *User `json:"user" valid:"null"`
	Token string `gorm:"unique_index" json:"-" valid:"null"`
	// 是否有效
	Effective bool `json:"effective" valid:"null"`
	// 登录方式
	LoginMethod LoginMethod `json:"login_method" valid:"null"`
	// 登录IP
	LoginIP   string     `gorm:"size:20" json:"login_ip"`
	CreatedAt time.Time  `json:"-" valid:"null"`
	UpdatedAt time.Time  `json:"-" valid:"null"`
	DeletedAt *time.Time `json:"-" sql:"index" valid:"null"`
}

// Equal 是否相等
func (session UserSession) Equal(sessionB UserSession) bool {
	if session.ID == sessionB.ID &&
		session.UserID == sessionB.UserID &&
		session.Effective == sessionB.Effective &&
		session.LoginMethod == sessionB.LoginMethod &&
		session.LoginIP == sessionB.LoginIP {
		return true
	}
	return false
}
