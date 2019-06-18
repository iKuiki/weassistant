package services

import (
	"fmt"
	"github.com/go-redis/redis"
	"weassistant/conf/rediskey"
	"weassistant/models"

	"github.com/jinzhu/gorm"
)

// AdministratorSessionService 管理员会话服务
type AdministratorSessionService interface {
	// 标准Mysql服务
	Get(sessionID uint64) (session models.AdministratorSession, err error)
	GetByWhereOptions(whereOptions []OrmWhereOption) (session models.AdministratorSession, err error)
	GetListByWhereOptions(whereOptions []OrmWhereOption, order []string, limit, offset int64, preloads ...string) (sessions []models.AdministratorSession, err error)
	GetCountByWhereOptions(whereOptions []OrmWhereOption) (count uint64, err error)
	Save(session *models.AdministratorSession) (err error)
	Delete(session *models.AdministratorSession) (err error)
	DeleteByWhereOptions(where []OrmWhereOption) (err error)
	// Redis附加服务
	ValidSessionToken(administratorID uint64, token string) (effective bool, err error)
}

type administratorSessionService struct {
	db            *gorm.DB
	commonService *commonService
	rds           *redis.Client
}

// MustNewAdministratorSessionService 新建管理员会话存储服务
func MustNewAdministratorSessionService(db *gorm.DB, rds *redis.Client) AdministratorSessionService {
	db.AutoMigrate(models.AdministratorSession{})
	return &administratorSessionService{
		db:            db,
		commonService: mustNewCommonService(db),
		rds:           rds,
	}
}

// Get 通过ID获取管理员会话
func (serv *administratorSessionService) Get(sessionID uint64) (session models.AdministratorSession, err error) {
	err = serv.commonService.Get(&session, sessionID)
	return
}

// GetByWhereOptions 通过查询条件获取管理员会话
func (serv *administratorSessionService) GetByWhereOptions(whereOptions []OrmWhereOption) (session models.AdministratorSession, err error) {
	err = serv.commonService.GetObjectByWhereOptions(&session, whereOptions)
	return
}

// GetListByWhereOptions 通过查询条件获取管理员会话列表
func (serv *administratorSessionService) GetListByWhereOptions(whereOptions []OrmWhereOption, order []string, limit, offset int64, preloads ...string) (sessions []models.AdministratorSession, err error) {
	err = serv.commonService.GetObjectListByWhereOptions(&sessions, whereOptions, order, limit, offset, preloads...)
	return
}

// GetCountByWhereOptions 通过查询条件获取管理员会话数量
func (serv *administratorSessionService) GetCountByWhereOptions(whereOptions []OrmWhereOption) (count uint64, err error) {
	return serv.commonService.GetCountByWhereOptions(models.AdministratorSession{}, whereOptions)
}

// Save 保存管理员会话
func (serv *administratorSessionService) Save(session *models.AdministratorSession) (err error) {
	if !session.Effective {
		err = serv.rds.SRem(fmt.Sprintf(rediskey.AdministratorTokenSet, session.AdministratorID), session.Token).Err()
	}
	if err != nil {
		return
	}
	err = serv.commonService.Save(session)
	if err == nil && session.Effective {
		err = serv.rds.SAdd(fmt.Sprintf(rediskey.AdministratorTokenSet, session.AdministratorID), session.Token).Err()
	}
	return
}

// Delete 删除管理员会话
func (serv *administratorSessionService) Delete(session *models.AdministratorSession) (err error) {
	s, err := serv.Get(session.ID)
	if err != nil {
		return
	}
	err = serv.rds.SRem(fmt.Sprintf(rediskey.AdministratorTokenSet, s.AdministratorID), s.Token).Err()
	if err != nil {
		return
	}
	return serv.commonService.Delete(session)
}

// DeleteByWhereOptions 根据查询条件删除管理员会话
func (serv *administratorSessionService) DeleteByWhereOptions(where []OrmWhereOption) (err error) {
	sessions, err := serv.GetListByWhereOptions(where, []string{}, 0, 0)
	if err != nil {
		return
	}
	idTokensMap := make(map[uint64][]interface{})
	for _, session := range sessions {
		idTokensMap[session.AdministratorID] = append(idTokensMap[session.AdministratorID], session.Token)
	}
	for id, tokens := range idTokensMap {
		err = serv.rds.SRem(fmt.Sprintf(rediskey.AdministratorTokenSet, id), tokens...).Err()
		if err != nil {
			return
		}
	}
	return serv.commonService.DeleteByWhereOptions(models.AdministratorSession{}, where)
}

// ------------------------ Redis附加功能 ------------------------
// 验证会话是否可用
func (serv *administratorSessionService) ValidSessionToken(administratorID uint64, token string) (effective bool, err error) {
	effective, err = serv.rds.SIsMember(fmt.Sprintf(rediskey.AdministratorTokenSet, administratorID), token).Result()
	return
}
