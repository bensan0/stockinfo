package models

type StockBasicInfo struct {
	Id          uint     `gorm:"unique;autoIncrement" json:"id"`
	Code        string   `gorm:"primarykey" json:"code"`
	Name        string   `json:"name"` //全名
	Industry    Industry `gorm:"foreignkey:IndustryId;references:Id;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;comment:業種" json:"industry"`
	IndustryId  uint     `gorm:"column:industry_id" json:"industry_id"`              //業種外鍵
	Capital     string   `gorm:"comment:資本額" json:"capital"`                         //資本額
	MarketValue string   `gorm:"column:market_value;comment:市值" json:"market_value"` //市值
	Intro       string   `json:"intro"`                                              //介紹
	UpdatedDate int      `json:"-"`
}

func (s StockBasicInfo) TableName() string {
	return "stock_basic_info"
}
