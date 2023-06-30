package examples

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func ProvinceGinBefore(c *gin.Context) {
	var err error
	if DBInit() == nil {
		c.String(http.StatusOK, GetProvinceBefore("3"))
	} else {
		c.String(http.StatusOK, MockGetProvinceBefore("3"))
	}
	if err != nil {
		println(err)
	}
}

func CityGinBefore(c *gin.Context) {
	var err error
	if DBInit() == nil {
		c.String(http.StatusOK, GetCityBefore("4"))
	} else {
		c.String(http.StatusOK, MockGetCityBefore("4"))
	}
	if err != nil {
		println(err)
	}
}

func ServerGinBefore() {
	r := gin.Default()
	r.GET("/province", ProvinceGinBefore)
	r.GET("/city", CityGinBefore)
	err := r.Run("127.0.0.1:2025")
	if err != nil {
		println(err)
	}
}
