package conf

import (
	"github.com/dgrijalva/jwt-go"
	"github.com/go-redis/redis"
	"github.com/jinzhu/gorm"
	"github.com/kataras/iris"
	"github.com/pkg/errors"
	"strings"
	"weassistant/conf/rediskey"
	"weassistant/services/locker"
	"weassistant/services/orm"

	jwtmiddleware "github.com/iris-contrib/middleware/jwt"
)

// ExtraConfig 扩展配置，基于基础配置，创建基于基础配置的服务抽象
type ExtraConfig interface {
	// GetJwtValidationKey 获取jwt认证key
	GetJwtValidationKey() (jwtValidationKey string)
	// GetAPIJwtMiddleware 获取前台api用jwt认证中间件
	GetAPIJwtMiddleware() (apiJwtMiddleware *jwtmiddleware.Middleware)
	// GetMgrJwtMiddleware 获取后台管理模块用jwt认证中间件
	GetMgrJwtMiddleware() (mgrJwtMiddleware *jwtmiddleware.Middleware)
	// GetRegisterLockerService 获取用户注册锁服务
	GetRegisterLockerService() (registerLockerService locker.CommonLockerService)
	// GetUserService 获取用户存取服务
	GetUserService() (userService orm.UserService)
	//  GetUserSessionService 获取用户Session存取服务
	GetUserSessionService() (userSessionService orm.UserSessionService)
	//  GetAdministratorService 获取管理员存取服务
	GetAdministratorService() (administratorService orm.AdministratorService)
	//  GetAdministratorSessionService 获取管理员Session存取服务
	GetAdministratorSessionService() (administratorSessionService orm.AdministratorSessionService)
}

// MustExtraNewConfig 必须创建扩展配置
func MustExtraNewConfig(c Config) ExtraConfig {
	e, err := NewExtraConfig(c)
	if err != nil {
		panic(err)
	}
	return e
}

// NewExtraConfig 根据给出的Config创建扩展配置
func NewExtraConfig(c Config) (extraC ExtraConfig, err error) {
	e := &extraConfig{
		jwtValidationKey: c.GetJwtValidationKey(),
	}
	err = errors.WithStack(e.RegisterJwtMiddleware())
	if err != nil {
		return
	}
	err = errors.WithStack(e.RegisterLockerService(c.GetMainRedis()))
	if err != nil {
		return
	}
	err = errors.WithStack(e.RegisterOrmService(c.GetMainDB(), c.GetMainRedis()))
	if err != nil {
		return
	}
	extraC = e
	return
}

// 扩展配置，基于基础配置，创建基于基础配置的服务
type extraConfig struct {
	// Jwt token相关
	jwtValidationKey string                    // Jwt认证key
	apiJwtMiddleware *jwtmiddleware.Middleware // 前台API的jwt验证中间件
	mgrJwtMiddleware *jwtmiddleware.Middleware // 后台管理模块的jwt验证中间件（与前台的区别在于jwt存放的header字段不同
	// 锁服务
	registerLockerService locker.CommonLockerService // 注册锁生成器
	// Orm服务
	userService                 orm.UserService                 // 用户存储服务
	userSessionService          orm.UserSessionService          // 用户Session存储服务
	administratorService        orm.AdministratorService        // 管理员存储服务
	administratorSessionService orm.AdministratorSessionService // 管理员Session存储服务
}

func (c *extraConfig) RegisterJwtMiddleware() (err error) {
	c.apiJwtMiddleware = jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(c.jwtValidationKey), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
		ErrorHandler: func(ctx iris.Context, err string) {
			ctx.Application().Logger().Debug("jwtHandler check error: ", err)
		},
	})
	// 创建后台管理模块使用的jwt认证中间件（与前台的不同在于jwt token存储的字段不同，后台是存在header的Authorization字段里
	c.mgrJwtMiddleware = jwtmiddleware.New(jwtmiddleware.Config{
		ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
			return []byte(c.jwtValidationKey), nil
		},
		SigningMethod: jwt.SigningMethodHS256,
		ErrorHandler: func(ctx iris.Context, err string) {
			ctx.Application().Logger().Debug("jwtHandler check error: ", err)
		},
		// 自定义提取器，从Authorization中提取jwt token
		Extractor: func(ctx iris.Context) (string, error) {
			authHeader := ctx.GetHeader("Authorization")
			if authHeader == "" {
				return "", nil // No error, just no token
			}
			authHeaderParts := strings.Split(authHeader, " ")
			if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearermgr" {
				return "", errors.Errorf("Authorization header format must be Bearer {token}")
			}
			return authHeaderParts[1], nil
		},
	})
	return nil
}

func (c *extraConfig) RegisterLockerService(redisClient *redis.Client) (err error) {
	c.registerLockerService, err = locker.NewCommonLockerService(redisClient, rediskey.UserRegisterLocker)
	if err != nil {
		return
	}
	return nil
}

func (c *extraConfig) RegisterOrmService(db *gorm.DB, redisClient *redis.Client) (err error) {
	c.userService, err = orm.NewUserService(db)
	if err != nil {
		return errors.WithStack(err)
	}
	c.userSessionService, err = orm.NewUserSessionService(db, redisClient)
	if err != nil {
		return errors.WithStack(err)
	}
	c.administratorService, err = orm.NewAdministratorService(db)
	if err != nil {
		return errors.WithStack(err)
	}
	c.administratorSessionService, err = orm.NewAdministratorSessionService(db, redisClient)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (c *extraConfig) GetJwtValidationKey() (jwtValidationKey string) {
	return c.jwtValidationKey
}

func (c *extraConfig) GetAPIJwtMiddleware() (apiJwtMiddleware *jwtmiddleware.Middleware) {
	return c.apiJwtMiddleware
}

func (c *extraConfig) GetMgrJwtMiddleware() (mgrJwtMiddleware *jwtmiddleware.Middleware) {
	return c.mgrJwtMiddleware
}

func (c *extraConfig) GetRegisterLockerService() (registerLockerService locker.CommonLockerService) {
	return c.registerLockerService
}

func (c *extraConfig) GetUserService() (userService orm.UserService) {
	return c.userService
}

func (c *extraConfig) GetUserSessionService() (userSessionService orm.UserSessionService) {
	return c.userSessionService
}

func (c *extraConfig) GetAdministratorService() (administratorService orm.AdministratorService) {
	return c.administratorService
}

func (c *extraConfig) GetAdministratorSessionService() (administratorSessionService orm.AdministratorSessionService) {
	return c.administratorSessionService
}
