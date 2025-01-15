package main

import (
	"fmt"
	"leetify-test/database"
	"leetify-test/routes"
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.SetConfigName("app.conf")

	err := viper.ReadInConfig()
	if err != nil {
		fmt.Println("Error reading config file:", err)
		return
	}
	var dbConf database.DBConnection
	err = viper.Sub("database").Unmarshal(&dbConf)
	if err != nil {
		fmt.Println("Error unmarshaling database config:", err)
		return
	}
	db, err := database.Connect(dbConf)
	if err != nil {
		fmt.Println("Error connecting to database:", err)
		return
	}
	database.Migrate(db)

	tmphttpreadheadertimeout, _ := time.ParseDuration(viper.GetString("server.readheadertimeout") + "s")
	tmphttpreadtimeout, _ := time.ParseDuration(viper.GetString("server.readtimeout") + "s")
	tmphttpwritetimeout, _ := time.ParseDuration(viper.GetString("server.writetimeout") + "s")
	tmphttpidletimeout, _ := time.ParseDuration(viper.GetString("server.idletimeout") + "s")

	r := gin.New()
	r.Use(cors.New(cors.Config{
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "User-Agent", "Referrer", "Host", "Token", "Authorization", "x-api-name", "x-api-key"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowAllOrigins:  false,
		AllowOriginFunc:  func(origin string) bool { return true },
	}))
	r.Use(gin.Recovery())

	routes.Routing(r, db)

	s := &http.Server{
		Addr:              ":" + viper.GetString("server.port"),
		Handler:           r,
		ReadHeaderTimeout: tmphttpreadheadertimeout,
		ReadTimeout:       tmphttpreadtimeout,
		WriteTimeout:      tmphttpwritetimeout,
		IdleTimeout:       tmphttpidletimeout,
	}

	fmt.Println("🚀 Server running on port:", viper.GetString("server.port"))
	s.ListenAndServe()
}
