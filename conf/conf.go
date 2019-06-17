package conf

import (
	"github.com/jinzhu/configor"
	"github.com/pkg/errors"
)

// Config 配置载体，所有可配置项都记载于此接口中
type Config interface {
	// Load 从文件中载入配置
	Load(filename string) error
}

// MustNewConfig 必须创建配置
func MustNewConfig() Config {
	c, err := NewConfig()
	if err != nil {
		panic(err)
	}
	return c
}

// NewConfig 创建配置
func NewConfig() (c Config, err error) {
	c = &config{}
	return
}

// config Config接口的默认实现
type config struct {
}

func (c *config) Load(filename string) error {
	return errors.WithStack(configor.Load(c, filename))
}
