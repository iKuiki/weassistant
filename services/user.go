package services

import (
	"weassistant/models"

	"github.com/jinzhu/gorm"
)

// UserService 用户服务
type UserService interface {
	Get(userID uint64) (user models.User, err error)
	GetByWhereOptions(whereOptions []OrmWhereOption) (user models.User, err error)
	GetListByWhereOptions(whereOptions []OrmWhereOption, order []string, limit, offset int64, preloads ...string) (users []models.User, err error)
	GetCountByWhereOptions(whereOptions []OrmWhereOption) (count uint64, err error)
	Save(user *models.User) (err error)
	Delete(user *models.User) (err error)
	DeleteByWhereOptions(where []OrmWhereOption) (err error)
}

type userService struct {
	db            *gorm.DB
	commonService *commonService
}

// MustNewUserService 新建用户存储服务
func MustNewUserService(db *gorm.DB) UserService {
	db.AutoMigrate(models.User{})
	return &userService{
		db:            db,
		commonService: mustNewCommonService(db),
	}
}

// Get 通过ID获取用户
func (serv *userService) Get(userID uint64) (user models.User, err error) {
	err = serv.commonService.Get(&user, userID)
	return
}

// GetByWhereOptions 通过查询条件获取用户
func (serv *userService) GetByWhereOptions(whereOptions []OrmWhereOption) (user models.User, err error) {
	err = serv.commonService.GetObjectByWhereOptions(&user, whereOptions)
	return
}

// GetListByWhereOptions 通过查询条件获取用户列表
func (serv *userService) GetListByWhereOptions(whereOptions []OrmWhereOption, order []string, limit, offset int64, preloads ...string) (users []models.User, err error) {
	err = serv.commonService.GetObjectListByWhereOptions(&users, whereOptions, order, limit, offset, preloads...)
	return
}

// GetCountByWhereOptions 通过查询条件获取用户数量
func (serv *userService) GetCountByWhereOptions(whereOptions []OrmWhereOption) (count uint64, err error) {
	return serv.commonService.GetCountByWhereOptions(models.User{}, whereOptions)
}

// Save 保存用户
func (serv *userService) Save(user *models.User) (err error) {
	return serv.commonService.Save(user)
}

// Delete 删除用户
func (serv *userService) Delete(user *models.User) (err error) {
	return serv.commonService.Delete(user)
}

// DeleteByWhereOptions 根据查询条件删除用户
func (serv *userService) DeleteByWhereOptions(where []OrmWhereOption) (err error) {
	return serv.commonService.DeleteByWhereOptions(models.User{}, where)
}
