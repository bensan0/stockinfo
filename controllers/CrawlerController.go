package controllers

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/bensan0/stockinfo/controllers/crawlers"
	models "github.com/bensan0/stockinfo/models/response"
)

type CrawlerController struct {
	web.Controller
}

func (this *CrawlerController) GetDailyQuot() {
	yyyymmdd := this.GetString("date")
	_, err := crawlers.DailyQuot(yyyymmdd)
	res := models.CommonRes{}
	if err != nil {
		res.Error = err.Error()
	}
	this.Data["json"] = &res
	this.ServeJSON()
}

func (this *CrawlerController) GetCorpDaily() {
	yyyymmdd := this.GetString("date")
	err := crawlers.DailyCorp(yyyymmdd)
	res := models.CommonRes{}
	if err != nil {
		res.Error = err.Error()
	}
	this.Data["json"] = &res
	this.ServeJSON()
}

func (this *CrawlerController) GetDistribution() {
	res := models.CommonRes{}
	err := crawlers.TWCCDCrawl()
	if err != nil {
		res.Error = err.Error()
	}
	this.Data["json"] = &res
	this.ServeJSON()
}
