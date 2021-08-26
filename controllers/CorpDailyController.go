package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/beego/beego/v2/server/web"
	common "github.com/bensan0/stockinfo/models/common"
	response "github.com/bensan0/stockinfo/models/response"
	"github.com/gomodule/redigo/redis"
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

	//快取
	if common.Cache != nil && days <= 15 {
		for _, date := range dates {
			mp := map[string]common.CorporationDailyTrans{}
			key := "corp" + fmt.Sprint(date)
			rpl, err := redis.String(common.Cache.Do("Get", key))
			if len(rpl) == 2 {
				continue
			}
			if err != nil {
				this.Data["json"] = &response.CommonRes{Error: err.Error()}
				this.ServeJSON()
			}
			json.Unmarshal([]byte(rpl), &mp)
			trans = append(trans, mp[code])
		}
		this.Data["json"] = &response.CommonRes{Data: trans}
		this.ServeJSON()
		return
	}

	var sql string
	if len(code) == 0 {
		sql = "date in ?"
	} else {
		sql = "date in ? and code = ?"
	}
	result := common.DB.Where(sql, dates, code).Order("date desc").Order("id").Find(&trans)

	if result.Error != nil {
		this.Data["json"] = &response.CommonRes{Error: result.Error.Error()}
		this.ServeJSON()
		return
	}

	this.Data["json"] = &response.CommonRes{Data: trans}
	this.ServeJSON()
}
