package easylog

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/span/field"
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
	attr := field.NewAttribute("123", nil)

	wg.Wait()

	l.Info("1")
	l.Info("2")
	l.Info("3")
	l.Info("4")
	l.Info("5")
	l.Info("6")

	l.Info("this  is a tst", field.WithAttribute(attr))
	l.Info("this  is a tst", field.WithAttribute(field.NewAttribute("123", field.StringField("fasfasfasf"))))
	type A struct {
		Name string `json:"name"`
		Age  string `json:"age"`
	}
	var a = A{
		Name: "zhangsan",
		Age:  "123",
	}
	l.InfoField(field.MallocJsonField(a), "test")
	l.Close()

}
