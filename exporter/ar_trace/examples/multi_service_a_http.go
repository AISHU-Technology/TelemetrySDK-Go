package examples

import (
	"io"
	"net/http"
)

func CheckAddressBefore() string {
	response, _ := http.Get("http://127.0.0.1:2023/province")
	province, _ := io.ReadAll(response.Body)

	response, _ = http.Get("http://127.0.0.1:2023/city")
	city, _ := io.ReadAll(response.Body)
	return " Address : " + string(province) + " Province " + string(city) + " City "
}
