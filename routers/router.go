package routers

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/bensan0/stockinfo/controllers"
)

func init() {
	web.Router("/", &controllers.MainController{})
	web.Router("/index", &controllers.MainController{}, "*:Index")

	ns := web.NewNamespace("/financial2",
		web.NSNamespace("/stockdaily",
			web.NSRouter("/getdays", &controllers.StockDailyController{}, "Get:GetDays"),
			web.NSRouter("/filttradingvol", &controllers.StockDailyController{}, "Get:FiltTradingVol"),
		),
		web.NSNamespace("/corpdaily",
			web.NSRouter("/getdays", &controllers.CorpDailyController{}, "Get:GetDays"),
		),
		web.NSNamespace("/crawlers",
			web.NSRouter("/dailyquot", &controllers.CrawlerController{}, "Get:GetDailyQuot"),
			web.NSRouter("/corpdaily", &controllers.CrawlerController{}, "Get:GetCorpDaily"),
			web.NSRouter("/distribution", &controllers.CrawlerController{}, "Get:GetDistribution"),
		),
		web.NSNamespace("/industry",
			web.NSRouter("/industry", &controllers.IndustryController{}),
		),
		web.NSNamespace("/stockinfo",
			web.NSRouter("/info", &controllers.StockBasicInfoController{}),
		),
	)
	web.AddNamespace(ns)
}
