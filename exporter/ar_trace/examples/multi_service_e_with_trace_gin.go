package examples

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"net/http"
)

func ProvinceGin(c *gin.Context) {
	var err error
	if DBInit() == nil {
		c.String(http.StatusOK, GetProvince(c.Request.Context(), "3"))
	} else {
		c.String(http.StatusOK, MockGetProvince(c.Request.Context(), "3"))
	}
	if err != nil {
		println(err)
	}
}

func CityGin(c *gin.Context) {
	var err error
	if DBInit() == nil {
		c.String(http.StatusOK, GetCity(c.Request.Context(), "4"))
	} else {
		c.String(http.StatusOK, MockGetCity(c.Request.Context(), "4"))
	}
	if err != nil {
		println(err)
	}
}

func ServerGin() {
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	r := gin.Default()
	r.Use(otelgin.Middleware("ServerGin"))

	r.GET("/province", ProvinceGin)
	r.GET("/city", CityGin)
	err := r.Run("127.0.0.1:2025")
	if err != nil {
		println(err)
	}
}
