package controllers

import (
	"github.com/beego/beego/v2/server/web"
	"github.com/bensan0/stockinfo/models/common"
)

type IndustryController struct {
	web.Controller
}

type res struct {
	Data  []models.Industry `json:"data"`
	Error error             `json:"error"`
}

func (this *IndustryController) Get() {
	industs := []models.Industry{}
	res := res{}
	models.DB.Find(&industs)
	res.Data = industs
	this.Data["json"] = res
	this.ServeJSON()
}

func (this *IndustryController) Post() {
	industs := []models.Industry{{}}
	res := res{}

	industs[0].Comment = this.GetString("comment")
	result := models.DB.Create(&industs[0])
	res.Data = industs
	res.Error = result.Error
	this.Data["json"] = res
	this.ServeJSON()
}

func (this *IndustryController) Delete() {
	res := res{}
	result := models.DB.Where("comment = ?", this.GetString("comment")).Delete(models.Industry{})
	res.Error = result.Error
	this.Data["json"] = res
	this.ServeJSON()
}
