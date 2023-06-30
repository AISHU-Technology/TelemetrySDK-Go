package examples

import (
	"context"
	"devops.aishu.cn/AISHUDevOps/ONE-Architecture/_git/TelemetrySDK-Go.git/exporter/ar_trace"
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/propagation"
	"io"
	"net/http"
)

func CheckAddress() string {
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	contextEmpty, span := ar_trace.Tracer.Start(context.Background(), "CheckAddress")
	defer span.End()

	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	bag, _ := baggage.Parse("username=donuts")
	ctx := baggage.ContextWithBaggage(contextEmpty, bag)

	var province, city []byte
	err := func(ctx context.Context) error {
		req, _ := http.NewRequestWithContext(ctx, "GET", "http://127.0.0.1:2023/province", nil)
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
		req, _ := http.NewRequestWithContext(ctx, "GET", "http://127.0.0.1:2023/city", nil)
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
	return " Address : " + string(province) + " Province " + string(city) + " City "
}
