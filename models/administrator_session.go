package models

import "time"

// AdministratorSession 管理员登录记录
type AdministratorSession struct {
	ID              uint64 `gorm:"primary_key" json:"id" form:"id" valid:"null"`
	AdministratorID uint64 `gorm:"index" json:"administrator_id" valid:"null"`
	// Administrator   *Administrator `json:"administrator" valid:"null"`
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
func (session AdministratorSession) Equal(sessionB AdministratorSession) bool {
	if session.ID == sessionB.ID &&
		session.AdministratorID == sessionB.AdministratorID &&
		session.Effective == sessionB.Effective &&
		session.LoginMethod == sessionB.LoginMethod &&
		session.LoginIP == sessionB.LoginIP {
		return true
	}
	return false
}
