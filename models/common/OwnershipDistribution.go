package models

//股權分佈

type OwnershipDistribution struct {
	Id   uint `gorm:"unique;autoincrement"`
	Code string
	//人數/張數/比例
	Rate1  string `gorm:"column:rate1;comment:持股<1(張) 人數/張數/比例"`
	Rate2  string `gorm:"column:rate2;comment:持股1-5(張)"`
	Rate3  string `gorm:"column:rate3;comment:持股5-10(張)"`
	Rate4  string `gorm:"column:rate4;comment:持股10-15(張)"`
	Rate5  string `gorm:"column:rate5;comment:持股15-20(張)"`
	Rate6  string `gorm:"column:rate6;comment:持股20-30(張)"`
	Rate7  string `gorm:"column:rate7;comment:持股30-40(張)"`
	Rate8  string `gorm:"column:rate8;comment:持股40-50(張)"`
	Rate9  string `gorm:"column:rate9;comment:持股50-100(張)"`
	Rate10 string `gorm:"column:rate10;comment:持股100-200(張)"`
	Rate11 string `gorm:"column:rate11;comment:持股200-400(張)"`
	Rate12 string `gorm:"column:rate12;comment:持股400-600(張)"`
	Rate13 string `gorm:"column:rate13;comment:持股600-800(張)"`
	Rate14 string `gorm:"column:rate14;comment:持股800-1000(張)"`
	Rate15 string `gorm:"column:rate15;comment:持股>1000(張)"`
	Total  string `gorm:"column:total;comment:r1-r15合計"`

	Date    int64  //yyyymmdd
	CDUnion string `gorm:"primarykey;column:cd_union;comment:代號-日期組合,做複合主鍵"` //複合主鍵
}

func (o OwnershipDistribution) TableName() string {
	return "ownership_distri"
}
