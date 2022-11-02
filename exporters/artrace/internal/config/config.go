package config

import (
	"time"
)

// Config Exporter的配置项
type Config struct {
	HTTPConfig  HTTPConfig
	RetryConfig RetryConfig
}

// DefaultConfig 默认的配置项。
var DefaultConfig = Config{
	HTTPConfig:  DefaultHTTPConfig,
	RetryConfig: DefaultRetryConfig,
}

// DefaultHTTPConfig 默认的HTTP配置。
var DefaultHTTPConfig = HTTPConfig{
	Insecure:    true,
	Endpoint:    "localhost:5678",
	Path:        "/api/feed_ingester/v1/jobs/traceTest/events",
	Compression: 1,
	Timeout:     10 * time.Second,
	Headers:     nil,
	TLSCfg:      nil,
}

// DefaultRetryConfig 默认的重发机制。
var DefaultRetryConfig = RetryConfig{
	Enabled:         true,
	InitialInterval: 5 * time.Second,
	MaxInterval:     30 * time.Second,
	MaxElapsedTime:  time.Minute,
}

// NewConfig 返回执行 Option 后的配置项。
func NewConfig(opts ...Option) Config {
	cfg := DefaultConfig
	for _, opt := range opts {
		cfg = opt.applyOption(cfg)
	}
	return cfg
}

// Option client的配置项结构体。
type Option struct {
	Fn func(Config) Config
}

func (h *Option) applyOption(cfg Config) Config {
	return h.Fn(cfg)
}
func newOption(fn func(cfg Config) Config) Option {
	return Option{Fn: fn}
}

func EmptyOption() Option {
	return newOption(func(cfg Config) Config {
		return cfg
	})
}
