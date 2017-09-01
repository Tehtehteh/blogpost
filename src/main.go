package main

import (
	format "fmt"
	"github.com/valyala/fasthttp"
	"github.com/buaazp/fasthttprouter"
	"api"
	"log"
)

func prepareRoutes(router *fasthttprouter.Router){
	router.GET("/", api.HomepageHanlder)
}

func main(){
	go format.Print("Serving...")
	router := fasthttprouter.New()
	prepareRoutes(router)
	log.Fatal(fasthttp.ListenAndServe("192.168.0.106:8090", router.Handler))
}