package models

//三大法人

type CorporationDailyTrans struct {
	Id               uint   `gorm:"unique;autoIncrement" json:"id"`
	Code             string `json:"code"`
	Name             string `json:"name"`
	ForeignInvestors int    `gorm:"column:foreign_investors;comment:外陸資總計買賣超" json:"foreignInvestors"` //外陸資：國外的投資機構或公司透過證劵經紀商或經紀部，在台灣購買股票。。
	ForeignCorp      int    `gorm:"column:foreign_corp;comment:外資自營商總計買賣超" json:"foreignCorp"`         //外資自營商:國外證劵商的自營部自行買賣股票
	InvestmentTrust  int    `gorm:"column:investment_trust;comment:投信總計買賣超" json:"investmentTrust"`    //投信：在國內發行基金，並使用這些資金進行投資的業者。
	Dealer           int    `gorm:"comment:自營商總計買賣超" json:"dealer"`                                    //自營商：使用自己公司資金做證券投資，又分為「專業自營商」或「證券公司自營部」。
	DealerSelf       int    `gorm:"column:dealer_self;comment:自營商(自行買賣)總計買賣超" json:"dealerSelf"`
	DealerHedge      int    `gorm:"column:dealer_hedge;comment:自營商(避險)總計買賣超" json:"dealerHedge"`
	Total            int    `gorm:"comment:三大法人合計買賣超(不計入自營商買賣)" json:"total"`                  //不計入Dealer
	Date             int64  `json:"date"`                                                      //yyyymmdd
	CDUnion          string `gorm:"primarykey;column:cd_union;comment:代號-日期組合,做複合主鍵" json:"-"` //複合主鍵
}

func (c CorporationDailyTrans) TableName() string {
	return "corporation_daily_trans"
}
