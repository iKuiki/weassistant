package conf

import (
	"fmt"
	"github.com/getsentry/raven-go"
	"github.com/go-redis/redis"
	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	// init mysql
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Config 配置载体，所有可配置项都记载于此接口中
type Config interface {
	// GetAppPort 获取程序运行端口
	GetAppPort() (port int)
	// GetDebug 是否为调试模式
	GetDebug() (isDebug bool)
	// Load 从文件中载入配置
	Load(filename string) error
	// GetMainDB 获取主DB
	GetMainDB() (dataDB *gorm.DB)
	// GetMainRedis 获取主Redis
	GetMainRedis() (client *redis.Client)
	// GetSentryClient 获取Sentry客户端
	GetSentryClient() (client *raven.Client)
	// GetJwtValidationKey 获取jwt认证key
	GetJwtValidationKey() (jwtValidationKey string)
}

// MustNewConfig 必须创建配置
func MustNewConfig() Config {
	c, err := NewConfig()
	if err != nil {
		panic(err)
	}
	return c
}

// 数据库连接的配置详情
type dbConf struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

// Redis连接的配置详情
type redisConf struct {
	Host     string
	Port     int
	Password string
	DBNo     int
	PoolSize int
}

// NewConfig 创建配置
func NewConfig() (c Config, err error) {
	c = &config{}
	return
}

// config Config接口的默认实现，包含从config中读取的配置以及根据其创建的连接
type config struct {
	// AppPort 运行端口
	AppPort int
	// 是否调试模式
	Debug bool
	// 数据库连接
	MainDBConf dbConf
	MainDB     *gorm.DB `json:"-"`
	// Redis连接
	MainRedisConf redisConf
	MainRedis     *redis.Client `json:"-"`
	// SentryDSN 错误捕获程序Sentry的路径
	SentryDSN    string
	SentryClient *raven.Client `json:"-"`
	// Jwt认证Key
	JwtValidationKey string
}

func (c *config) Load(filename string) (err error) {
	err = errors.WithStack(configor.Load(c, filename))
	if err != nil {
		return err
	}
	err = errors.WithStack(c.RegisterMainDB())
	if err != nil {
		return err
	}
	err = errors.WithStack(c.RegisterMainRedis())
	if err != nil {
		return err
	}
	err = errors.WithStack(c.RegisterSentry())
	if err != nil {
		return err
	}
	return nil
}

func (c *config) RegisterMainDB() (err error) {
	c.MainDB, err = gorm.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=True&loc=Local",
			c.MainDBConf.User,
			c.MainDBConf.Password,
			c.MainDBConf.Host,
			c.MainDBConf.Port,
			c.MainDBConf.DBName,
		))
	return
}

func (c *config) RegisterMainRedis() (err error) {
	mainRedis := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", c.MainRedisConf.Host, c.MainRedisConf.Port),
		Password: c.MainRedisConf.Password,
		DB:       c.MainRedisConf.DBNo,
		PoolSize: c.MainRedisConf.PoolSize,
	})
	if err := mainRedis.Ping().Err(); err != nil {
		return err
	}
	c.MainRedis = mainRedis
	return nil
}

func (c *config) RegisterSentry() (err error) {
	if c.SentryDSN != "" {
		// 如果SentryDSN为空，则认为不启用Sentry
		sentryClient, err := raven.New(c.SentryDSN)
		if err != nil {
			return errors.WithStack(err)
		}
		c.SentryClient = sentryClient
	}
	return nil
}

func (c *config) GetAppPort() (port int) {
	return c.AppPort
}

func (c *config) GetDebug() (isDebug bool) {
	return c.Debug
}

func (c *config) GetMainDB() (dataDB *gorm.DB) {
	return c.MainDB
}

func (c *config) GetMainRedis() (client *redis.Client) {
	return c.MainRedis
}

func (c *config) GetSentryClient() (client *raven.Client) {
	return c.SentryClient
}

func (c *config) GetJwtValidationKey() (jwtValidationKey string) {
	return c.JwtValidationKey
}
