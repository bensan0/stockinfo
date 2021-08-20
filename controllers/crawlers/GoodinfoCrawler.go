package crawlers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/beego/beego/v2/server/web"
	common "github.com/bensan0/stockinfo/models/common"
	"github.com/bensan0/stockinfo/tools"
	"github.com/chromedp/chromedp"
	"gorm.io/gorm"
)

//從goodinfo爬個股資料

func StockBasicInfoCrawl(code string) error {
	sbi, err := DownloadStockBasicInfo(code)
	if err != nil {
		tools.ErrorLogging(err)
		return err
	}
	err = InfoUpsert(sbi)
	if err != nil {
		tools.ErrorLogging(err)
	}
	return err
}

//個股基本資料
func DownloadStockBasicInfo(code string) (*common.StockBasicInfo, error) {
	fmt.Println("stockBasicInfoCrawl 開始!!!!!!!!!!!!!")
	url, _ := web.AppConfig.String("stockbasicurl")
	url += code
	//1配置options
	//1-1.自定義設定
	options := []chromedp.ExecAllocatorOption{
		chromedp.Flag("headless", false),                       // debug使用
		chromedp.Flag("blink-settings", "imagesEnabled=false"), //不加載圖片
		chromedp.UserAgent(`Mozilla/5.0 (Windows NT 6.3; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.103 Safari/537.36`),
	}
	//1-2.以自定義設定取代原本的預設設定
	options = append(chromedp.DefaultExecAllocatorOptions[:], options...)

	//1-3.取得新設定
	alloCtx, cancel := chromedp.NewExecAllocator(context.Background(), options...)
	defer cancel()

	//2.配置context
	//2-1. create chrome instance
	ctx, cancel := chromedp.NewContext(
		alloCtx,
		chromedp.WithLogf(log.Printf),
	)
	defer cancel()
	//2-2. create a timeout
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	//3.訪問目標網頁並執行動作
	var basicInfohtml string
	err := chromedp.Run(
		ctx,
		chromedp.Navigate(url),

		//獲取公司基本資訊
		chromedp.OuterHTML(`body > table:nth-child(8) > tbody > tr > td:nth-child(3) > table > tbody > tr:nth-child(2) > td:nth-child(3) > table:nth-child(1)`, &basicInfohtml),
	)
	if err != nil {
		return nil, err
	}
	basicInfohtml = "<html><head></head>" + basicInfohtml + "</html>"
	fmt.Println("stockBasicInfoCrawl 結束!!!!!!!!!!!!!")
	basicInfoptr, err := basicInfoFilter(code, basicInfohtml)
	if err != nil {
		return nil, err
	}
	return basicInfoptr, nil
}

//解析
func basicInfoFilter(stockCode, basicInfohtml string) (*common.StockBasicInfo, error) {
	basicInfo := common.StockBasicInfo{}
	dom, _ := goquery.NewDocumentFromReader(strings.NewReader(basicInfohtml))
	basicInfo.Code = stockCode
	basicInfo.UpdatedDate, _ = strconv.Atoi(time.Now().Format("20060102"))
	//資本額
	dom.Find(`body > table > tbody > tr:nth-child(5) > td:nth-child(2) > nobr`).Each(func(i int, s *goquery.Selection) {
		basicInfo.Capital = s.Text()
	})

	if len(basicInfo.Capital) == 0 {
		return &basicInfo, errors.New("無內容")
	}
	//市值
	dom.Find(`body > table > tbody > tr:nth-child(5) > td:nth-child(4) > nobr`).Each(func(i int, s *goquery.Selection) {
		basicInfo.MarketValue = s.Text()
	})
	//介紹
	dom.Find(`body > table > tbody > tr:nth-child(15) > td > p`).Each(func(i int, s *goquery.Selection) {
		basicInfo.Intro = s.Text()
	})
	//全名
	dom.Find(`body > table > tbody > tr:nth-child(2) > td:nth-child(2)`).Each(func(i int, s *goquery.Selection) {
		basicInfo.Name = s.Text()
	})
	//業種
	dom.Find(`body > table > tbody > tr:nth-child(3) > td:nth-child(2)`).Each(func(i int, s *goquery.Selection) {
		indust := s.Text()

		if !tools.Contains(common.IndustCommentList, indust) {
			temp := &common.Industry{Comment: indust}
			common.DB.Create(temp)
			basicInfo.IndustryId = temp.Id
		} else {
			basicInfo.IndustryId = common.IndustCommentMap[indust]
		}
	})
	return &basicInfo, nil
}

//Update or Insert
func InfoUpsert(info *common.StockBasicInfo) error {
	old := common.StockBasicInfo{}
	common.DB.Where("code = ?", info.Code).Find(&old)
	var result *gorm.DB
	if old.Id == 0 {
		result = common.DB.Create(info)
	} else {
		result = common.DB.Select("*").Omit("Id").Updates(*info)
	}
	if result.Error != nil {
		return result.Error
	}
	return nil
}
