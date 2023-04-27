package config

type KafkaConfig struct {
	Address  []string `json:"address" validate:"required"`
	User     string   `json:"user" validate:"required"`
	Password string   `json:"password" validate:"required"`
	Topic    string   `json:"topic" validate:"required"`
}

// 以下各项为发送器配置项目。

// WithAnyRobotURL 设置上报地址。
func WithAddress(address []string) Option {
	return func(cfg *Config) *Config {
		cfg.KafkaConfig.Address = address
		return cfg
	}
}

// WithUser 设置sasl认证用户名。
func WithUser(user string) Option {
	return func(cfg *Config) *Config {
		cfg.KafkaConfig.User = user
		return cfg
	}
}

// WithPassword 设置sasl认证密码。
func WithPassword(password string) Option {
	return func(cfg *Config) *Config {
		cfg.KafkaConfig.Password = password
		return cfg
	}
}

// WithTopic 设置主题。
func WithTopic(topic string) Option {
	return func(cfg *Config) *Config {
		cfg.KafkaConfig.Topic = topic
		return cfg
	}
}
