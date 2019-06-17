package conf

// Config 配置载体，所有可配置项都记载于此接口中
type Config interface {
}

// config Config接口的默认实现
type config struct {
}
