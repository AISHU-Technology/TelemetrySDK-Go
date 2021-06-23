package benchmarks

import (
	"gitlab.aishu.cn/anyrobot/observability/telemetrysdk/telemetry-go/span/field"
	"testing"
)

func BenchmarkGenSpanID(b *testing.B) {
	b.Log("BenchmarkGenSpanID")

	b.Run("field/GenSpanID", func(b *testing.B) {
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				field.GenSpanID()
			}
		})
	})
}
