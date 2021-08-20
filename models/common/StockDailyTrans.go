package models

type StockDailyTrans struct {
	Id           uint    `gorm:"unique;autoIncrement" json:"id"`
	Code         string  `json:"code"`
	Name         string  `json:"name"`
	TradingVol   int     `gorm:"column:trading_vol;comment:成交量(張)" json:"tradingVol"`       //成交量(張)
	Deal         int     `gorm:"comment:成交筆數" json:"deal"`                                  //成交筆數
	Opening      float64 `gorm:"comment:開盤價" json:"opening"`                                //開盤
	Closing      float64 `gorm:"comment:收盤價" json:"closing"`                                //收盤
	Highest      float64 `gorm:"comment:當日最高" json:"highest"`                               //當日最高
	Lowest       float64 `gorm:"comment:當日最低" json:"lowest"`                                //當日最低
	Fluctuation  float64 `gorm:"comment:漲跌價差" json:"fluctuation"`                           //漲跌價差 = 今收與昨收價差
	FluctPercent float64 `gorm:"column:fluc_percent;comment:漲跌幅(%)" json:"flucPercent"`     //漲跌幅(%)
	Date         int64   `json:"date"`                                                      //yyyymmdd
	CDUnion      string  `gorm:"primarykey;column:cd_union;comment:代號-日期組合,做複合主鍵" json:"-"` //複合主鍵
}

func (s StockDailyTrans) TableName() string {
	return "stock_daily_trans"
}
