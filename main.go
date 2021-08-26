package main

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/filter/cors"
	common "github.com/bensan0/stockinfo/models/common"

	_ "github.com/bensan0/stockinfo/routers"
)

func init() {
	common.DBInit()
	enableRedis, _ := web.AppConfig.Bool("enableRedis")
	if enableRedis{
		common.CachInit()
	}
	initCors()
}

func main() {
	web.Run()
}

func initCors() {
	web.InsertFilter("*", web.BeforeRouter, cors.Allow(&cors.Options{
		AllowAllOrigins: true,
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Authorization", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		ExposeHeaders: []string{"Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers", "Content-Type"},
		AllowCredentials: true,
	}))
}
