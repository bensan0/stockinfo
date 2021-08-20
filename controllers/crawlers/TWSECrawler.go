package crawlers

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/server/web"
	common "github.com/bensan0/stockinfo/models/common"
	"golang.org/x/text/encoding/traditionalchinese"
	"golang.org/x/text/transform"
)

//證交所

//每日三大法人調用
func DailyCorp(yyyymmdd string) error {
	filename, err := DownloadDailyCorpTrans(yyyymmdd)
	if err != nil {
		return err
	}
	err = DailyCorpFilter(filename)
	if err != nil {
		return err
	}
	return nil
}

//下載三大法人每日交易csv
func DownloadDailyCorpTrans(yyyymmdd string) (string, error) {
	fmt.Println("開始下載每日三大法人")
	now := time.Now()
	if weekday := now.Weekday(); weekday == time.Saturday || weekday == time.Sunday {
		return "", errors.New("六/日休市")
	}
	if len(yyyymmdd) == 0 {
		yyyymmdd = now.Format("20060102")
	}
	var filedir string = "downloads/"
	var filename string = "dailycorp_"
	var sub string = ".csv"
	url, _ := web.AppConfig.String("corpdailyurl")
	url = url + yyyymmdd
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filedir + filename + yyyymmdd + sub)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Write the body to file
	wr, err := io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	} else if wr < 10 {
		return "", errors.New("無內容")
	}
	fmt.Println("每日三大法人下載結束")
	return filedir + filename + yyyymmdd + sub, nil
}

//解析
func DailyCorpFilter(filename string) error {
	corpdailys := []common.CorporationDailyTrans{}

	fmt.Println("每日三大法人處理開始")
	by, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	r, _, _ := transform.String(traditionalchinese.Big5.NewDecoder(), string(by))
	arr := strings.Split(r, `"三大法人買賣超股數",`)
	arr = strings.Split(arr[1], `"說明:"`)
	sc := bufio.NewScanner(strings.NewReader(strings.TrimSpace(arr[0])))
	for sc.Scan() {
		str := string(sc.Bytes())
		str = str[1 : len(str)-2]           //去除頭尾符號以便分割
		strArr := strings.Split(str, `","`) //分割

		//數據轉換為模型
		corpdaily := common.CorporationDailyTrans{}

		corpdaily.Code = strings.Replace(strArr[0], `"`, "", 1)

		corpdaily.Name = strings.TrimSpace(strArr[1])

		fitemp, _ := strconv.Atoi(strings.Replace(strArr[4], ",", "", -1))
		corpdaily.ForeignInvestors = fitemp / 1000

		fctemp, _ := strconv.Atoi(strings.Replace(strArr[7], ",", "", -1))
		corpdaily.ForeignCorp = fctemp / 1000

		ittemp, _ := strconv.Atoi(strings.Replace(strArr[10], ",", "", -1))
		corpdaily.InvestmentTrust = ittemp / 1000

		dtemp, _ := strconv.Atoi(strings.Replace(strArr[11], ",", "", -1))
		corpdaily.Dealer = dtemp / 1000

		dstemp, _ := strconv.Atoi(strings.Replace(strArr[14], ",", "", -1))
		corpdaily.DealerSelf = dstemp / 1000

		dhtemp, _ := strconv.Atoi(strings.Replace(strArr[17], ",", "", -1))
		corpdaily.DealerHedge = dhtemp / 1000

		corpdaily.Date, _ = strconv.ParseInt(filename[20:28], 10, 64)

		corpdaily.CDUnion = corpdaily.Code + "-" + fmt.Sprint(corpdaily.Date)

		ttemp, _ := strconv.Atoi(strings.Replace(strArr[18], ",", "", -1))
		corpdaily.Total = ttemp / 1000
		corpdailys = append(corpdailys, corpdaily)
	}
	err = os.Remove(filename)
	if err != nil {
		return err
	}
	fmt.Println("每日三大法人處理結束")

	fmt.Println("存入DB開始")
	result := common.DB.Create(corpdailys)
	if result.Error != nil {
		fmt.Println("存入DB失敗")
		return result.Error
	}
	fmt.Println("存入DB結束")
	return nil
}

//每日收盤行情調用
func DailyQuot(yyyymmdd string) (bool, error) {
	filename, err := DownloadDailyQuotation(yyyymmdd)
	if err != nil || len(filename) == 0 {
		return false, err
	}
	err = DailyQuotFilter(filename)
	if err != nil {
		return false, err
	}
	return true, err
}

//下載每日收盤行情(全部(不含權證、牛熊證))csv
func DownloadDailyQuotation(yyyymmdd string) (string, error) {
	fmt.Println("開始下載每日收盤行情")

	//日期處理
	now := time.Now()
	if weekday := now.Weekday(); weekday == time.Saturday || weekday == time.Sunday {
		return "", errors.New("六/日休市")
	}
	if len(yyyymmdd) == 0 {
		yyyymmdd = now.Format("20060102")
	}
	var filedir string = "downloads/"
	var filename string = "dailyquot_"
	var sub string = ".csv"
	url, _ := web.AppConfig.String("stockdailyurl")
	url = url + yyyymmdd
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filedir + filename + yyyymmdd + sub)
	if err != nil {
		return "", err
	}
	defer out.Close()

	// Write the body to file
	wr, err := io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	} else if wr == 0 {
		return "", err
	}
	fmt.Println("每日收盤行情下載結束")
	return filedir + filename + yyyymmdd + sub, nil
}

//處理檔案每日收盤行情檔案
func DailyQuotFilter(filename string) error {
	fmt.Println("每日收盤行情處理開始")
	stdailys := []common.StockDailyTrans{}
	by, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil
	}

	//字元編碼
	r, _, _ := transform.String(traditionalchinese.Big5.NewDecoder(), string(by))
	arr := strings.Split(r, `"本益比",`)

	//循環處理
	sc := bufio.NewScanner(strings.NewReader(strings.TrimSpace(arr[1])))
	for sc.Scan() {
		str := string(sc.Bytes())
		str = str[1 : len(str)-2]           //去除尾部符號以便分割
		strArr := strings.Split(str, `","`) //分割

		//數據轉換為模型
		stdaily := common.StockDailyTrans{}
		stdaily.Date, _ = strconv.ParseInt(filename[20:28], 10, 64)
		strArr[0] = strings.Replace(strArr[0], `"`, "", 1)
		stdaily.Code = strArr[0]
		stdaily.Name = strArr[1]
		tv, _ := strconv.Atoi(strings.Replace(strArr[2], ",", "", -1))
		stdaily.TradingVol = tv / 1000
		stdaily.Deal, _ = strconv.Atoi(strArr[3])
		stdaily.Opening, _ = strconv.ParseFloat(strArr[5], 64)
		stdaily.Highest, _ = strconv.ParseFloat(strArr[6], 64)
		stdaily.Lowest, _ = strconv.ParseFloat(strArr[7], 64)
		stdaily.Closing, _ = strconv.ParseFloat(strArr[8], 64)

		if strArr[9] == "-" {
			fluTemp, _ := strconv.ParseFloat(strArr[10], 64)
			stdaily.Fluctuation = -fluTemp
		} else {
			stdaily.Fluctuation, _ = strconv.ParseFloat(strArr[10], 64)
		}
		stdaily.FluctPercent, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", stdaily.Fluctuation/(stdaily.Closing-stdaily.Fluctuation)*100), 64)
		stdaily.CDUnion = stdaily.Code + "-" + fmt.Sprint(stdaily.Date)
		stdailys = append(stdailys, stdaily)

	}

	//刪檔
	err = os.Remove(filename)
	if err != nil {
		return err
	}
	fmt.Println("每日收盤行情處理結束")

	//持久化
	fmt.Println("存入DB開始")
	result := common.DB.Create(stdailys)
	if result.Error != nil {
		fmt.Println("存入DB失敗")
		return result.Error
	}
	fmt.Println("存入DB結束")
	return nil
}
