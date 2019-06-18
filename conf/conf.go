package conf

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/jinzhu/configor"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"

	// init mysql
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

// Config 配置载体，所有可配置项都记载于此接口中
type Config interface {
	// Load 从文件中载入配置
	Load(filename string) error
	// GetMainDB 获取主DB
	GetMainDB() (dataDB *gorm.DB)
	// GetMainRedis 获取主Redis
	GetMainRedis() (client *redis.Client)
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

// config Config接口的默认实现
type config struct {
	MainDBConf    dbConf
	MainDB        *gorm.DB `json:"-"`
	MainRedisConf redisConf
	MainRedis     *redis.Client `json:"-"`
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

func (c *config) GetMainDB() (dataDB *gorm.DB) {
	return c.MainDB
}

func (c *config) GetMainRedis() (client *redis.Client) {
	return c.MainRedis
}
