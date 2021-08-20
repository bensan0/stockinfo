package models

//連接MySqlDB

import (
	"fmt"
	"strings"

	"github.com/beego/beego/v2/server/web"
	"github.com/bensan0/stockinfo/tools"

	_ "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var err error
var dialect string = "mysqladmin:mysqlpassword@/mysqldb?charset=utf8&parseTime=true&loc=Local"
var IndustIdMap = map[uint]string{}      //[id]comment
var StockCodeMap = map[string]int{}      //[code]industryid
var IndustCommentList = []string{}       //[comment]
var IndustCommentMap = map[string]uint{} //[comment]id

func DBInit() {
	connectDB()
	err := DB.AutoMigrate(
		&StockBasicInfo{},
		&Industry{},
		&StockDailyTrans{},
		&CorporationDailyTrans{},
		&OwnershipDistribution{},
	)
	if err != nil {
		tools.ErrorLogging(err)
	}
	getIndustryAndStock()
}

func connectDB() {
	mysqladmin, _ := web.AppConfig.String("mysqladmin")
	mysqlpassword, _ := web.AppConfig.String("mysqlpassword")
	mysqldb, _ := web.AppConfig.String("mysqldb")

	address := strings.Replace(dialect, "mysqladmin", mysqladmin, -1)
	address = strings.Replace(address, "mysqlpassword", mysqlpassword, -1)
	address = strings.Replace(address, "mysqldb", mysqldb, -1)

	DB, err = gorm.Open(mysql.Open(address), &gorm.Config{Logger: logger.Default.LogMode(logger.Info)})
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("DB連接成功")
	}
}

func getIndustryAndStock() {
	industs := []Industry{}
	stocks := []StockBasicInfo{}
	DB.Find(&industs)
	DB.Find(&stocks)
	for _, v := range industs {
		IndustIdMap[v.Id] = v.Comment
		IndustCommentList = append(IndustCommentList, v.Comment)
		IndustCommentMap[v.Comment] = v.Id
	}
	for _, v := range stocks {
		StockCodeMap[v.Code] = int(v.IndustryId)
	}
}
