package config

import (
	"context"
	customErrors "devops.aishu.cn/AISHUDevOps/AnyRobot/_git/Akashic_TelemetrySDK-Go.git/exporters/artrace/internal/errors"
	"errors"
	"github.com/cenkalti/backoff/v4"
	"log"
	"time"
)

// RetryConfig 当Trace数据发送失败时，根据重发机制来重新发送，保障数据不漏。
type RetryConfig struct {
	// Enabled 是否启用重发机制。
	Enabled bool
	// InitialInterval 第一次重发与上一次发送的时间间隔。
	InitialInterval time.Duration
	// MaxInterval 两次重发的最长时间间隔。
	MaxInterval time.Duration
	// MaxElapsedTime 重发最长持续的时间。
	MaxElapsedTime time.Duration
}

// RetryFunc 包含重发机制的请求方法。
type RetryFunc func(context.Context, func(context.Context) error) error

// RetryableError 可以触发重发机制的错误。
type RetryableError struct {
	Throttle int64
}

// Error 实现error接口。
func (e RetryableError) Error() string {
	return customErrors.AnyRobotTraceExporter_RetryFailure
}

// evaluate 通过 RetryableError 类型来判断是否可重发。
func evaluate(err error) (bool, time.Duration) {
	if err == nil {
		return false, 0
	}

	rErr, ok := err.(RetryableError)
	if !ok {
		return false, 0
	}

	return true, time.Duration(rErr.Throttle)
}

// RetryFunc 带有重发机制的HTTP请求，如果错误为可重发错误则重发AnyRobotSpans，否则丢弃数据。
func (c RetryConfig) RetryFunc() RetryFunc {
	if !c.Enabled {
		return func(ctx context.Context, fn func(context.Context) error) error {
			return fn(ctx)
		}
	}

	b := &backoff.ExponentialBackOff{
		InitialInterval:     c.InitialInterval,
		RandomizationFactor: backoff.DefaultRandomizationFactor,
		Multiplier:          backoff.DefaultMultiplier,
		MaxInterval:         c.MaxInterval,
		MaxElapsedTime:      c.MaxElapsedTime,
		Stop:                backoff.Stop,
		Clock:               backoff.SystemClock,
	}
	b.Reset()

	//返回一个嵌套的函数是为了重发时复用HTTP连接，外部的 retryFunc(context,function)error 是包装了retry的逻辑，内部的function(context),error是真正的http请求逻辑，在client处实现。
	retryFunc := func(ctx context.Context, fn func(context.Context) error) error {
		for {
			err := fn(ctx)
			if err == nil {
				return nil
			}

			retryable, throttle := evaluate(err)
			if !retryable {
				return err
			}

			backOff := b.NextBackOff()
			if backOff == backoff.Stop {
				log.Println(customErrors.AnyRobotTraceExporter_ExceedRetryElapsedTime)
				return errors.New(customErrors.AnyRobotTraceExporter_ExceedRetryElapsedTime)
			}

			var delay time.Duration
			if backOff > throttle {
				delay = backOff
			} else {
				elapsed := b.GetElapsedTime()
				if b.MaxElapsedTime != 0 && elapsed+throttle > b.MaxElapsedTime {
					log.Println(customErrors.AnyRobotTraceExporter_ExceedRetryElapsedTime)
					return errors.New(customErrors.AnyRobotTraceExporter_ExceedRetryElapsedTime)
				}
				delay = throttle
			}

			if err := wait(ctx, delay); err != nil {
				return err
			}
		}
	}
	return retryFunc
}

// wait 等待一段时间。
func wait(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		select {
		case <-timer.C:
		default:
			return ctx.Err()
		}
	case <-timer.C:
	}

	return nil
}

// WithRetry 设置重发。
func WithRetry(rc RetryConfig) HTTPOption {
	return newHTTPOption(func(cfg Config) Config {
		cfg.RetryConfig = rc
		return cfg
	})
}
