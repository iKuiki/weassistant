package models

import (
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// User 用户
type User struct {
	ID uint64 `gorm:"primary_key" json:"id"`
	// Nickname 名字
	Nickname string `gorm:"size:60" json:"nickname" form:"nickname" validate:"required,min=2,max=20"`
	// Account 登陆名
	Account string `gorm:"unique_index;size:60" json:"account" form:"account" validate:"required,min=4,max=20"`
	// Password bcrypt加密后的密码
	Password string `gorm:"size:140" json:"-" form:"password" validate:"printascii,min=6,max=40"`
	// 最后登陆时间
	LastLoginAt *time.Time `json:"last_login_at"`
	// 最后登陆IP
	LastLoginIP string `json:"last_login_ip"`
	// 素质三连
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"-"`
}

// CreateTo 根据表单输入创建到数据对象上
func (user User) CreateTo(dest *User) {
	dest.Nickname = user.Nickname
	dest.Account = user.Account
	dest.Password = user.Password
}

// UpdateTo 根据输入更新到数据对象上
func (user User) UpdateTo(dest *User) {
	dest.Nickname = user.Nickname
	if user.Password != "" {
		dest.Password = user.Password
	}
}

// Equal 是否相等
func (user User) Equal(userB User) bool {
	if user.ID == userB.ID &&
		user.Nickname == userB.Nickname &&
		user.Account == userB.Account &&
		// user.LastLoginAt == userB.LastLoginAt &&
		user.LastLoginIP == userB.LastLoginIP {
		return true
	}
	return false
}

// BeforeSave 保存前预处理
func (user *User) BeforeSave() (err error) {
	// 判断用户密码未加密则用bcrypt进行加密（依据为不以$2a$04$开头且长度小于或等于40
	if !strings.HasPrefix(user.Password, "$2a$04$") && len(user.Password) <= 40 {
		hashedPass, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.MinCost)
		if err != nil {
			return err
		}
		user.Password = string(hashedPass)
	}
	return nil
}
