package controllers

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/server/web"
	common "github.com/bensan0/stockinfo/models/common"
	response "github.com/bensan0/stockinfo/models/response"
	"github.com/gomodule/redigo/redis"
)

type stDailyRes struct {
	Code     string      `json:"code"`
	Name     string      `json:"name"`
	Industry string      `json:"industry"`
	Trans    interface{} `json:"trans"`
}

type StockDailyController struct {
	web.Controller
}

//依據天數與代號獲取個股交易資訊
func (this *StockDailyController) GetDays() {

	trans := []common.StockDailyTrans{}
	today := time.Now()
	dates := make([]int64, 0)
	days, _ := this.GetInt("days", 5)
	code := strings.TrimSpace(this.GetString("code"))

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
			mp := map[string]common.StockDailyTrans{}
			key := "quot" + fmt.Sprint(date)
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

//篩選交易量>前日個股
//輸入:日期,百分比,交易量,
func (this *StockDailyController) FiltTradingVol() {

	dateStr := strings.TrimSpace(this.GetString("date"))
	if len(dateStr) == 0 {
		dateStr = time.Now().Format("20060102")
	}
	date, _ := time.Parse("20060102", dateStr)

	percent, _ := this.GetFloat("percent", 40)
	lowerLimitVol, _ := this.GetInt("lowervol", 3000)
	higherLimitVol, _ := this.GetInt("highervol", 100000)
	arr := []int64{}

	for i := 0; ; i++ {
		if len(arr) >= 2 {
			break
		}
		tmp := date.AddDate(0, 0, -i)
		if tmp.Weekday() == time.Saturday || tmp.Weekday() == time.Sunday {
			continue
		}
		before, _ := strconv.ParseInt(tmp.Format("20060102"), 10, 64)
		arr = append(arr, before)
	}

	todayMap := map[string]common.StockDailyTrans{}
	yesterdayMap := map[string]common.StockDailyTrans{}
	dailytrans := []common.StockDailyTrans{}
	result := common.DB.Where("date in ?", arr).Order("id").Find(&dailytrans)
	if result.Error != nil {
		this.Data["json"] = &response.CommonRes{Data: nil, Error: result.Error.Error()}
		this.ServeJSON()
		return
	}
	for _, dt := range dailytrans {
		if dt.Date == arr[0] {
			todayMap[dt.Code] = dt
		} else {
			yesterdayMap[dt.Code] = dt
		}
	}

	resultList := []stDailyRes{}
	percent = 1 + percent/100
	for k, v := range todayMap {
		tv := float64(v.TradingVol)
		ytv := float64(yesterdayMap[k].TradingVol)

		if len(v.Code) == 4 && int(ytv) >= lowerLimitVol && int(ytv) <= higherLimitVol && tv >= ytv*percent {
			resultList = append(resultList, stDailyRes{
				Code:     v.Code,
				Industry: common.IndustIdMap[uint(common.StockCodeMap[v.Code])],
				Name:     v.Name,
				Trans:    []common.StockDailyTrans{v, yesterdayMap[k]},
			})
		}
	}
	//依業種排序
	sort.Slice(resultList, func(i, j int) bool {
		return common.IndustCommentMap[resultList[i].Industry] < common.IndustCommentMap[resultList[j].Industry]
	})
	this.Data["json"] = &response.CommonRes{Data: resultList}
	this.ServeJSON()
}
