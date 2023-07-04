package examples

import (
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"net/http"
)

func Province(w http.ResponseWriter, r *http.Request) {
	var err error
	if DBInit() == nil {
		_, err = w.Write([]byte(GetProvince(r.Context(), "3")))
	} else {
		_, err = w.Write([]byte(MockGetProvince(r.Context(), "3")))
	}
	if err != nil {
		println(err)
	}
}

func City(w http.ResponseWriter, r *http.Request) {
	var err error
	if DBInit() == nil {
		_, err = w.Write([]byte(GetCity(r.Context(), "4")))
	} else {
		_, err = w.Write([]byte(MockGetCity(r.Context(), "4")))
	}
	if err != nil {
		println(err)
	}
}

func Server() {
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	http.Handle("/province", otelhttp.NewHandler(http.HandlerFunc(Province), "/province"))
	http.Handle("/city", otelhttp.NewHandler(http.HandlerFunc(City), "/city"))
	err := http.ListenAndServe("127.0.0.1:2023", nil)
	if err != nil {
		println(err)
	}

}
