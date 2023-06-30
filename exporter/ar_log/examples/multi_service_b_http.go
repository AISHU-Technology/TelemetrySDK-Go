package examples

import "net/http"

func ProvinceBefore(w http.ResponseWriter, r *http.Request) {
	var err error
	if DBInitBefore() == nil {
		_, err = w.Write([]byte(GetProvinceBefore("3")))
	} else {
		_, err = w.Write([]byte(MockGetProvinceBefore("3")))
	}
	if err != nil {
		println(err)
	}
}

func CityBefore(w http.ResponseWriter, r *http.Request) {
	var err error
	if DBInitBefore() == nil {
		_, err = w.Write([]byte(GetCityBefore("4")))
	} else {
		_, err = w.Write([]byte(MockGetCityBefore("4")))
	}
	if err != nil {
		println(err)
	}
}

func ServerBefore() {
	http.Handle("/province", http.HandlerFunc(ProvinceBefore))
	http.Handle("/city", http.HandlerFunc(CityBefore))
	err := http.ListenAndServe("127.0.0.1:2023", nil)
	if err != nil {
		println(err)
	}
}
