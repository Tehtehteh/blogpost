package main

import (
	"api"
	"db"
	"settings"
	. "logger"
	"github.com/valyala/fasthttp"
	"github.com/buaazp/fasthttprouter"
	"github.com/spf13/viper"
	"os"
	"fmt"
)

func prepareRoutes(router *fasthttprouter.Router){
	router.GET("/", api.HomepageHandler)
	router.POST("/b/:buyer_id", api.BuyerHandler)
	router.POST("/openrtb/", api.BuyerHandler)
	router.PanicHandler = api.PanicHandler
}

func main(){
	InitLog()
	Log.Info("Starting RTB...")
	var (
		configPath string
		connectionString string
	)
	if os.Getenv("DEBUG") == "1" {
		configPath += "config_develop"
	} else {
		configPath += "config"
	}
	viper.AddConfigPath("D:\\Golangprojects\\blogpost\\src\\config")
	viper.SetConfigName(configPath)
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()

	Log.Info("Qwa!")
	if err != nil {
		Log.Errorf("Error reading configuration file: %s", err)
		return
	}

	connectionString = fmt.Sprintf("%s:%s@%s%s", viper.GetString("mysql.user"), viper.GetString("mysql.password"),
		viper.GetString("mysql.host"), viper.GetString("mysql.database"))
	err = db.InitializeDB(connectionString)
	if err != nil {
		Log.Errorf("Error connecting to database: %s", err)
		return
	}
	settings.ExchangeHandler.Init()

	router := fasthttprouter.New()
	prepareRoutes(router)
	Log.Info(fasthttp.ListenAndServe(":8090", router.Handler))
}