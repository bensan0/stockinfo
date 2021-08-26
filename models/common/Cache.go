package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/bensan0/stockinfo/tools"
	"github.com/gomodule/redigo/redis"
)

var Cache redis.Conn

func CachInit() {
	Cache, err = redis.Dial("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("Redis連接失敗")
		tools.ErrorLogging(err)
	} else {
		fmt.Println("Redis連接成功")
	}
	cacheData()
}

func cacheData() {
	now := time.Now()
	dates := []string{}
	stdailys := []StockDailyTrans{}
	cpdailys := []CorporationDailyTrans{}

	for i := 0; ; i++ {
		if len(dates) >= 15 {
			break
		}

		date := now.AddDate(0, 0, -i)

		if date.Weekday() == time.Saturday || date.Weekday() == time.Sunday {
			continue
		}

		dates = append(dates, date.Format("20060102"))
	}

	DB.Where("date in ?", dates).Find(&stdailys)
	DB.Where("date in ?", dates).Find(&cpdailys)
	for i := 0; i < 15; i++ {
		quotkey := "quot" + dates[i]
		sts := map[string]StockDailyTrans{}
		for _, v := range stdailys {
			if fmt.Sprint(v.Date) == dates[i] {
				sts[v.Code] = v
			}
		}
		by, _ := json.Marshal(sts)
		Cache.Send("Set", quotkey, string(by))

		corpkey := "corp" + dates[i]
		cps := map[string]CorporationDailyTrans{}
		for _, v := range cpdailys {
			if fmt.Sprint(v.Date) == dates[i] {
				cps[v.Code] = v
			}
		}
		by, _ = json.Marshal(cps)
		Cache.Send("Set", corpkey, string(by))
	}
	Cache.Flush()
	if _, err := Cache.Receive();err!=nil{
		tools.ErrorLogging(err)
	}else{
		fmt.Println("Redis初始化成功")
	}
	
}
