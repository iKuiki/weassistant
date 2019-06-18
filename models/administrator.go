package models

import (
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

// Administrator 后台账号
type Administrator struct {
	ID uint64 `gorm:"primary_key" json:"id" validate:"isdefault"`
	// Name 名字
	Name string `gorm:"size:60" json:"name" form:"name" validate:"required,max=20"`
	// Account 登陆名
	Account string `gorm:"unique_index;size:60" json:"account" form:"account" validate:"required,min=4,max=20"`
	// Password bcrypt加密后的密码
	Password string `gorm:"size:140" json:"-" form:"password" validate:"printascii,min=6,max=40"`
	// 最后登陆时间
	LastLoginAt *time.Time `json:"last_login_at" validate:"isdefault"`
	// 最后登陆IP
	LastLoginIP string `json:"last_login_ip" validate:"isdefault"`
	// 素质三连
	CreatedAt time.Time  `json:"created_at" validate:"isdefault"`
	UpdatedAt time.Time  `json:"updated_at" validate:"isdefault"`
	DeletedAt *time.Time `sql:"index" json:"-" validate:"isdefault"`
}

// CreateTo 根据表单输入创建到数据对象上
func (administrator Administrator) CreateTo(dest *Administrator) {
	dest.Name = administrator.Name
	dest.Account = administrator.Account
	dest.Password = administrator.Password
}

// UpdateTo 根据输入更新到数据对象上
func (administrator Administrator) UpdateTo(dest *Administrator) {
	dest.Name = administrator.Name
	if administrator.Password != "" {
		dest.Password = administrator.Password
	}
}

// Equal 是否相等
func (administrator Administrator) Equal(administratorB Administrator) bool {
	if administrator.ID == administratorB.ID &&
		administrator.Name == administratorB.Name &&
		administrator.Account == administratorB.Account &&
		// administrator.LastLoginAt == administratorB.LastLoginAt &&
		administrator.LastLoginIP == administratorB.LastLoginIP {
		return true
	}
	return false
}

// BeforeSave 保存前预处理
func (administrator *Administrator) BeforeSave() (err error) {
	// 判断用户密码未加密则用bcrypt进行加密（依据为不以$2a$04$开头且长度小于或等于40
	if !strings.HasPrefix(administrator.Password, "$2a$04$") && len(administrator.Password) <= 40 {
		hashedPass, err := bcrypt.GenerateFromPassword([]byte(administrator.Password), bcrypt.MinCost)
		if err != nil {
			return err
		}
		administrator.Password = string(hashedPass)
	}
	return nil
}
