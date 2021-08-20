package controllers

import (
	"strconv"
	"time"

	"github.com/beego/beego/v2/server/web"
	common "github.com/bensan0/stockinfo/models/common"
	response "github.com/bensan0/stockinfo/models/response"
)

type CorpDailyController struct {
	web.Controller
}

func (this *CorpDailyController) GetDays() {
	code := this.GetString("code")
	days, _ := this.GetInt("days", 5)

	trans := []common.CorporationDailyTrans{}
	today := time.Now()
	dates := make([]int64, 0)

	for i := 0; ; i++ {
		if len(dates) >= days {
			break
		}
		tmp := today.AddDate(0, 0, -i)
		if tmp.Weekday() == time.Saturday || tmp.Weekday() == time.Sunday {
			continue
		}
		num, _ := strconv.ParseInt(tmp.Format("20060102"), 10, 64)
		dates = append(dates, num)
	}
	result := common.DB.Where("date in ? and code = ?", dates, code).Order("date desc").Order("id").Find(&trans)

	if result.Error != nil {
		this.Data["json"] = &response.CommonRes{Error: result.Error.Error()}
		this.ServeJSON()
		return
	}

	this.Data["json"] = &response.CommonRes{Data: trans}
	this.ServeJSON()
}
