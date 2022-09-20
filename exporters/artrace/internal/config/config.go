package config

import (
	"time"
)

type Config struct {
	HTTPConfig  HTTPConfig
	RetryConfig RetryConfig
}

// DefaultHTTPConfig 默认的重发机制。
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

// NewConfig 返回执行后的HTTP配置项。
func NewConfig(opts ...HTTPOption) Config {
	cfg := Config{
		HTTPConfig:  DefaultHTTPConfig,
		RetryConfig: DefaultRetryConfig,
	}
	for _, opt := range opts {
		cfg = opt.ApplyHTTPOption(cfg)
	}
	return cfg
}
