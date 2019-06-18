package services

import (
	"fmt"
	"github.com/go-redis/redis"
	"weassistant/conf/rediskey"
	"weassistant/models"

	"github.com/jinzhu/gorm"
)

// UserSessionService 用户会话服务
type UserSessionService interface {
	// 标准Mysql服务
	Get(sessionID uint64) (session models.UserSession, err error)
	GetByWhereOptions(whereOptions []OrmWhereOption) (session models.UserSession, err error)
	GetListByWhereOptions(whereOptions []OrmWhereOption, order []string, limit, offset int64, preloads ...string) (sessions []models.UserSession, err error)
	GetCountByWhereOptions(whereOptions []OrmWhereOption) (count uint64, err error)
	Save(session *models.UserSession) (err error)
	Delete(session *models.UserSession) (err error)
	DeleteByWhereOptions(where []OrmWhereOption) (err error)
	// Redis附加服务
	ValidSessionToken(userID uint64, token string) (effective bool, err error)
}

type userSessionService struct {
	db            *gorm.DB
	commonService *commonService
	rds           *redis.Client
}

// MustNewUserSessionService 新建用户会话存储服务
func MustNewUserSessionService(db *gorm.DB, rds *redis.Client) UserSessionService {
	db.AutoMigrate(models.UserSession{})
	return &userSessionService{
		db:            db,
		commonService: mustNewCommonService(db),
		rds:           rds,
	}
}

// Get 通过ID获取用户会话
func (serv *userSessionService) Get(sessionID uint64) (session models.UserSession, err error) {
	err = serv.commonService.Get(&session, sessionID)
	return
}

// GetByWhereOptions 通过查询条件获取用户会话
func (serv *userSessionService) GetByWhereOptions(whereOptions []OrmWhereOption) (session models.UserSession, err error) {
	err = serv.commonService.GetObjectByWhereOptions(&session, whereOptions)
	return
}

// GetListByWhereOptions 通过查询条件获取用户会话列表
func (serv *userSessionService) GetListByWhereOptions(whereOptions []OrmWhereOption, order []string, limit, offset int64, preloads ...string) (sessions []models.UserSession, err error) {
	err = serv.commonService.GetObjectListByWhereOptions(&sessions, whereOptions, order, limit, offset, preloads...)
	return
}

// GetCountByWhereOptions 通过查询条件获取用户会话数量
func (serv *userSessionService) GetCountByWhereOptions(whereOptions []OrmWhereOption) (count uint64, err error) {
	return serv.commonService.GetCountByWhereOptions(models.UserSession{}, whereOptions)
}

// Save 保存用户会话
func (serv *userSessionService) Save(session *models.UserSession) (err error) {
	if !session.Effective {
		err = serv.rds.SRem(fmt.Sprintf(rediskey.UserTokenSet, session.UserID), session.Token).Err()
	}
	if err != nil {
		return
	}
	err = serv.commonService.Save(session)
	if err == nil && session.Effective {
		err = serv.rds.SAdd(fmt.Sprintf(rediskey.UserTokenSet, session.UserID), session.Token).Err()
	}
	return
}

// Delete 删除用户会话
func (serv *userSessionService) Delete(session *models.UserSession) (err error) {
	s, err := serv.Get(session.ID)
	if err != nil {
		return
	}
	err = serv.rds.SRem(fmt.Sprintf(rediskey.UserTokenSet, s.UserID), s.Token).Err()
	if err != nil {
		return
	}
	return serv.commonService.Delete(session)
}

// DeleteByWhereOptions 根据查询条件删除用户会话
func (serv *userSessionService) DeleteByWhereOptions(where []OrmWhereOption) (err error) {
	sessions, err := serv.GetListByWhereOptions(where, []string{}, 0, 0)
	if err != nil {
		return
	}
	idTokensMap := make(map[uint64][]interface{})
	for _, session := range sessions {
		idTokensMap[session.UserID] = append(idTokensMap[session.UserID], session.Token)
	}
	for id, tokens := range idTokensMap {
		err = serv.rds.SRem(fmt.Sprintf(rediskey.UserTokenSet, id), tokens...).Err()
		if err != nil {
			return
		}
	}
	return serv.commonService.DeleteByWhereOptions(models.UserSession{}, where)
}

// ------------------------ Redis附加功能 ------------------------
func (serv *userSessionService) ValidSessionToken(userID uint64, token string) (effective bool, err error) {
	effective, err = serv.rds.SIsMember(fmt.Sprintf(rediskey.UserTokenSet, userID), token).Result()
	return
}
