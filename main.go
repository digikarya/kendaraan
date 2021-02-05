package main

import (
	"github.com/digikarya/kendaraan/app"
	"github.com/gorilla/mux"
	"log"
	"github.com/digikarya/kendaraan/config"
	"net/http"
)

// ServeHTTP wraps the HTTP server enabling CORS headers.
// For more info about CORS, visit https://www.w3.org/TR/cors/

func main() {
	r := mux.NewRouter()
	conf := config.GetConfig()
	//authHelper := &authHelper.Auth{}
	//authHelper.Initialize(conf,r)
	Kepegawaian := &app.Kendaraan{}
	Kepegawaian.Initialize(conf,r)
	//saksi := &app.SaksiApp{}
	//saksi.Initialize(config,r)
	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Print("App is not running")
	}else{
		log.Print("App is not running")
	}
}

