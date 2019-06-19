package orm

import (
	"github.com/jinzhu/gorm"
)

// WhereOption 传入where查询时使用的结构体
type WhereOption struct {
	Query string
	Item  []interface{}
}

type commonService struct {
	db *gorm.DB
}

// mustNewCommonService 创建通用存储服务
func mustNewCommonService(db *gorm.DB) *commonService {
	return &commonService{
		db: db,
	}
}

// ========================== Query ==========================

// Get 通过ID获取对象
func (serv *commonService) Get(object interface{}, id uint64) (err error) {
	err = serv.db.Where("id = ?", id).First(object).Error
	return
}

// GetObjectByWhereOptions 通过指定条件获取对象
func (serv *commonService) GetObjectByWhereOptions(object interface{}, where []WhereOption, preloads ...string) (err error) {
	o := serv.db
	for _, w := range where {
		o = o.Where(w.Query, w.Item...)
	}
	for _, preload := range preloads {
		o = o.Preload(preload)
	}
	err = o.First(object).Error
	return
}

// GetObjectListByWhereOptions 通过指定条件获取列表
func (serv *commonService) GetObjectListByWhereOptions(objects interface{}, where []WhereOption, orders []string, limit, offset int64, preloads ...string) (err error) {
	o := serv.db
	for _, w := range where {
		o = o.Where(w.Query, w.Item...)
	}
	for _, preload := range preloads {
		o = o.Preload(preload)
	}
	if len(orders) == 0 {
		o = o.Order("id")
	} else {
		for _, order := range orders {
			o = o.Order(order)
		}
	}
	if limit > 0 {
		o = o.Limit(limit)
	}
	if offset > 0 {
		o = o.Offset(offset)
	}
	// 此操作会执行赋值复制
	err = o.Find(objects).Error
	return
}

// GetCountByWhereOptions 通过指定条件获取计数
func (serv *commonService) GetCountByWhereOptions(model interface{}, where []WhereOption) (count uint64, err error) {
	o := serv.db.Model(model)
	for _, w := range where {
		o = o.Where(w.Query, w.Item...)
	}
	err = o.Count(&count).Error
	return
}

// ========================== Write ==========================

// Create 创建对象
func (serv *commonService) Create(object interface{}) (err error) {
	return serv.db.Create(object).Error
}

// Save 保存对象
func (serv *commonService) Save(object interface{}) (err error) {
	return serv.db.Save(object).Error
}

// UpdateObjectsViaMap 通过map更新对象
func (serv *commonService) UpdateObjectsViaMap(model interface{}, where []WhereOption, updateMap map[string]interface{}) (err error) {
	o := serv.db.Model(model)
	for _, w := range where {
		o = o.Where(w.Query, w.Item...)
	}
	err = o.Updates(updateMap).Error
	return
}

// ========================== Delete ==========================

// Delete 删除对象
func (serv *commonService) Delete(object interface{}) (err error) {
	return serv.db.Delete(object).Error
}

// DeleteByWhereOptions 根据条件删除对象
func (serv *commonService) DeleteByWhereOptions(model interface{}, where []WhereOption) (err error) {
	o := serv.db
	for _, w := range where {
		o = o.Where(w.Query, w.Item...)
	}
	return o.Delete(model).Error
}
