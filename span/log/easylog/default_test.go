package easylog

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/field"
)

func TestNewDefaultSamplerLogger(t *testing.T) {

	l := NewDefaultSamplerLogger()
	l.Info("this is a test", nil)
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 3; i++ {
			l.WarnField(field.StringField(fmt.Sprintf("%d", i)), "test", nil)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 3; i++ {
			l.WarnField(field.StringField(fmt.Sprintf("%d", i)), "test", nil)
		}
		time.Sleep(time.Second)
	}()

	l.Error("error", nil)
	attr := &field.Attribute{Message: field.StringField("123"), Type: "tsaga"}

	wg.Wait()

	l.Info("this  is a tst", attr)
	l.Close()

}
