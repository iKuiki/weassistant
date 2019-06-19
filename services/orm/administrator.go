package orm

import (
	"github.com/pkg/errors"
	"weassistant/models"

	"github.com/jinzhu/gorm"
)

// AdministratorService 管理员服务
type AdministratorService interface {
	Get(administratorID uint64) (administrator models.Administrator, err error)
	GetByWhereOptions(whereOptions []WhereOption) (administrator models.Administrator, err error)
	GetListByWhereOptions(whereOptions []WhereOption, order []string, limit, offset int64, preloads ...string) (administrators []models.Administrator, err error)
	GetCountByWhereOptions(whereOptions []WhereOption) (count uint64, err error)
	Save(administrator *models.Administrator) (err error)
	Delete(administrator *models.Administrator) (err error)
	DeleteByWhereOptions(where []WhereOption) (err error)
}

type administratorService struct {
	db            *gorm.DB
	commonService *commonService
}

// MustNewAdministratorService 新建管理员存储服务
func MustNewAdministratorService(db *gorm.DB) AdministratorService {
	serv, err := NewAdministratorService(db)
	if err != nil {
		panic(errors.WithStack(err))
	}
	return serv
}

// NewAdministratorService 新建管理员存储服务
func NewAdministratorService(db *gorm.DB) (serv AdministratorService, err error) {
	err = errors.WithStack(db.AutoMigrate(models.Administrator{}).Error)
	if err != nil {
		return
	}
	serv = &administratorService{
		db:            db,
		commonService: mustNewCommonService(db),
	}
	return
}

// Get 通过ID获取管理员
func (serv *administratorService) Get(administratorID uint64) (administrator models.Administrator, err error) {
	err = serv.commonService.Get(&administrator, administratorID)
	return
}

// GetByWhereOptions 通过查询条件获取管理员
func (serv *administratorService) GetByWhereOptions(whereOptions []WhereOption) (administrator models.Administrator, err error) {
	err = serv.commonService.GetObjectByWhereOptions(&administrator, whereOptions)
	return
}

// GetListByWhereOptions 通过查询条件获取管理员列表
func (serv *administratorService) GetListByWhereOptions(whereOptions []WhereOption, order []string, limit, offset int64, preloads ...string) (administrators []models.Administrator, err error) {
	err = serv.commonService.GetObjectListByWhereOptions(&administrators, whereOptions, order, limit, offset, preloads...)
	return
}

// GetCountByWhereOptions 通过查询条件获取管理员数量
func (serv *administratorService) GetCountByWhereOptions(whereOptions []WhereOption) (count uint64, err error) {
	return serv.commonService.GetCountByWhereOptions(models.Administrator{}, whereOptions)
}

// Save 保存管理员
func (serv *administratorService) Save(administrator *models.Administrator) (err error) {
	return serv.commonService.Save(administrator)
}

// Delete 删除管理员
func (serv *administratorService) Delete(administrator *models.Administrator) (err error) {
	return serv.commonService.Delete(administrator)
}

// DeleteByWhereOptions 根据查询条件删除管理员
func (serv *administratorService) DeleteByWhereOptions(where []WhereOption) (err error) {
	return serv.commonService.DeleteByWhereOptions(models.Administrator{}, where)
}
