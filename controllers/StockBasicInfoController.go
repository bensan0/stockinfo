package controllers

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/bensan0/stockinfo/controllers/crawlers"
	common "github.com/bensan0/stockinfo/models/common"
	response "github.com/bensan0/stockinfo/models/response"
)

type StockBasicInfoController struct {
	web.Controller
}

//get all
func (this *StockBasicInfoController) Get() {
	infos := []common.StockBasicInfo{{}}
	code := this.GetString("code")
	common.DB.Preload("Industry").Where("code = ?", code).Find(&infos)
	this.Data["json"] = infos
	this.ServeJSON()
}

func (this *StockBasicInfoController) Put() {
	code := this.GetString("code")
	res := response.CommonRes{}
	if len(code) == 0 {
		res.Error = "無code"
	} else {
		err := crawlers.StockBasicInfoCrawl(code)
		if err!=nil{
			res.Error = err.Error()
		}
	}
	this.Data["json"] = &res
	this.ServeJSON()
}
