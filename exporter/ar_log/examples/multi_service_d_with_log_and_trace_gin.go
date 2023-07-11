package examples

import (
	"context"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/ar_trace"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/propagation"
	"io"
	"net/http"
)

func CheckAddressGin() {
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	r := gin.Default()
	r.Use(otelgin.Middleware("CheckAddressGin"))

	r.GET("/address", func(c *gin.Context) {
		contextEmpty, span := ar_trace.Tracer.Start(context.Background(), "CheckAddress")
		client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
		bag, _ := baggage.Parse("username=donuts")
		ctx := baggage.ContextWithBaggage(contextEmpty, bag)

		SystemLogger.Warn("this is a test warn message")
		_ = ServiceLogger.Error("this is a test error message")

		var province, city []byte
		err := func(ctx context.Context) error {
			req, _ := http.NewRequestWithContext(ctx, "GET", "http://127.0.0.1:2025/province", nil)
			fmt.Println("Sending request...")
			res, err := client.Do(req)
			if err != nil {
				return err
			}
			province, err = io.ReadAll(res.Body)
			_ = res.Body.Close()
			return err
		}(ctx)
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("Response Received")
		}

		Err := func(ctx context.Context) error {
			req, _ := http.NewRequestWithContext(ctx, "GET", "http://127.0.0.1:2025/city", nil)
			fmt.Println("Sending request...")
			res, Err := client.Do(req)
			if Err != nil {
				return Err
			}
			city, Err = io.ReadAll(res.Body)
			_ = res.Body.Close()
			return Err
		}(ctx)
		if Err != nil {
			fmt.Println(Err)
		} else {
			fmt.Println("Response Received")
		}

		c.String(http.StatusOK, " Address : "+string(province)+" Province "+string(city)+" City ")
		span.End()
	})
	err := r.Run("127.0.0.1:2024")
	if err != nil {
		println(err)
	}
}
