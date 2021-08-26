package controllers

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/bensan0/stockinfo/controllers/crawlers"
	common "github.com/bensan0/stockinfo/models/common"
	response "github.com/bensan0/stockinfo/models/response"
)

type CrawlerController struct {
	web.Controller
}

func (this *CrawlerController) GetDailyQuot() {
	yyyymmdd := this.GetString("date")
	ch := make(chan []common.StockDailyTrans)
	res := response.CommonRes{}
	trans := []common.StockDailyTrans{}

	go func(yyyymmdd string) {
		quots, err := crawlers.DailyQuotTWSE(yyyymmdd)
		if err != nil {
			res.Error = err.Error()
			ch <- nil
		} else {
			ch <- quots
		}
	}(yyyymmdd)

	go func(yyyymmdd string) {
		quots, err := crawlers.DailyQuotPEX(yyyymmdd)
		if err != nil {
			res.Error = err.Error()
			ch <- nil
		} else {
			ch <- quots
		}
	}(yyyymmdd)

	for i := 0; i < 2; i++ {
		trans = append(trans, <-ch...)
	}

	err := crawlers.InsertToDB(trans)
	if err != nil {
		res.Error = err.Error()
	}

	this.Data["json"] = &res
	this.ServeJSON()
}

func (this *CrawlerController) GetCorpDaily() {
	yyyymmdd := this.GetString("date")
	ch := make(chan []common.CorporationDailyTrans)
	res := response.CommonRes{}
	trans := []common.CorporationDailyTrans{}

	go func(yyyymmdd string) {
		corps, err := crawlers.DailyCorpTPEX(yyyymmdd)
		if err != nil {
			res.Error = err.Error()
			ch <- nil
		} else {
			ch <- corps
		}
	}(yyyymmdd)

	go func(yyyymmdd string) {
		corps, err := crawlers.DailyCorpTWSE(yyyymmdd)
		if err != nil {
			res.Error = err.Error()
			ch <- nil
		} else {
			ch <- corps
		}
	}(yyyymmdd)

	for i := 0; i < 2; i++ {
		trans = append(trans, <-ch...)
	}

	err := crawlers.InsertToDB(trans)
	if err != nil {
		res.Error = err.Error()
	}
	this.Data["json"] = &res
	this.ServeJSON()
}

func (this *CrawlerController) GetDistribution() {
	res := response.CommonRes{}
	err := crawlers.TWCCDCrawl()
	if err != nil {
		res.Error = err.Error()
	}
	this.Data["json"] = &res
	this.ServeJSON()
}
