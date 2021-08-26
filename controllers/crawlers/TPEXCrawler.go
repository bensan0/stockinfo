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

//櫃買中心

func DailyCorpTPEX(yyyymmdd string) ([]common.CorporationDailyTrans, error) {
	filename, err := downloadDailyCorpTransTPEX(yyyymmdd)
	if err != nil {
		return nil, err
	}
	dailycorpTrans, err := dailyCorpFilterTPEX(filename)
	if err != nil {
		return nil, err
	}
	return dailycorpTrans, nil
}

func DailyQuotPEX(yyyymmdd string) ([]common.StockDailyTrans, error) {
	filename, err := downloadDailyQuotTPEX(yyyymmdd)
	if err != nil || len(filename) == 0 {
		return nil, err
	}
	stdailyTrans, err := dailyQuotFilterTPEX(filename)
	if err != nil {
		return nil, err
	}
	return stdailyTrans, err
}

func downloadDailyQuotTPEX(yyyymmdd string) (string, error) {
	fmt.Println("每日上櫃盤後下載開始")
	var yyyymmddTime time.Time
	if len(yyyymmdd) == 0 {
		yyyymmddTime = time.Now()
		yyyymmdd = yyyymmddTime.Format("20060102")

	} else {
		yyyymmddTime, _ = time.Parse("20060102", yyyymmdd)
	}
	if weekday := yyyymmddTime.Weekday(); weekday == time.Saturday || weekday == time.Sunday {
		return "", errors.New("六/日休市")
	}
	var filedir string = "downloads/"
	var filename string = "dailyquottpex_"
	var sub string = ".csv"
	url, _ := web.AppConfig.String("tpexdailyurl")
	url += yyyymmddTime.AddDate(-1911, 0, 0).Format("2006/01/02")[1:]
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
	} else if wr <= 96 {
		return "", errors.New("無內容")
	}
	fmt.Println("每日上櫃盤後下載結束")
	return filedir + filename + yyyymmdd + sub, nil
}

func dailyQuotFilterTPEX(filename string) ([]common.StockDailyTrans, error) {
	fmt.Println("每日上櫃盤後處理開始")
	stdailys := []common.StockDailyTrans{}
	by, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	//字元編碼
	r, _, _ := transform.String(traditionalchinese.Big5.NewDecoder(), string(by))
	arr := strings.Split(r, ",次日跌停價")
	arr = strings.Split(arr[1], "共")

	//循環處理
	sc := bufio.NewScanner(strings.NewReader(strings.TrimSpace(arr[0])))
	for sc.Scan() {
		str := string(sc.Bytes())
		str = str[1 : len(str)-1]           //去除尾部符號以便分割
		strArr := strings.Split(str, `","`) //分割

		//數據轉換為模型
		stdaily := common.StockDailyTrans{}
		stdaily.Date, _ = strconv.ParseInt(filename[24:32], 10, 64)
		stdaily.Code = strArr[0]
		stdaily.Name = strArr[1]

		clTemp := strings.Replace(strArr[2], ",", "", -1)
		stdaily.Closing, _ = strconv.ParseFloat(clTemp, 64)

		stdaily.Fluctuation, _ = strconv.ParseFloat(strArr[3], 64)

		opTemp := strings.Replace(strArr[4], ",", "", -1)
		stdaily.Opening, _ = strconv.ParseFloat(opTemp, 64)

		hiTemp := strings.Replace(strArr[5], ",", "", -1)
		stdaily.Highest, _ = strconv.ParseFloat(hiTemp, 64)

		loTemp := strings.Replace(strArr[6], ",", "", -1)
		stdaily.Lowest, _ = strconv.ParseFloat(loTemp, 64)

		tv, _ := strconv.Atoi(strings.Replace(strArr[7], ",", "", -1))
		stdaily.TradingVol = tv / 1000

		stdaily.Deal, _ = strconv.Atoi(strings.Replace(strArr[9], ",", "", -1))
		stdaily.FluctPercent, _ = strconv.ParseFloat(fmt.Sprintf("%.3f", stdaily.Fluctuation/(stdaily.Closing-stdaily.Fluctuation)*100), 64)
		stdaily.CDUnion = stdaily.Code + "-" + fmt.Sprint(stdaily.Date)
		stdailys = append(stdailys, stdaily)
	}
	// 刪檔
	err = os.Remove(filename)
	if err != nil {
		return nil, err
	}
	fmt.Println("每日上櫃盤後處理結束")

	return stdailys, nil
}

func downloadDailyCorpTransTPEX(yyyymmdd string) (string, error) {
	fmt.Println("每日上櫃三大法人下載開始")
	var yyyymmddTime time.Time
	if len(yyyymmdd) == 0 {
		yyyymmddTime = time.Now()
		yyyymmdd = yyyymmddTime.Format("20060102")

	} else {
		yyyymmddTime, _ = time.Parse("20060102", yyyymmdd)
	}
	if weekday := yyyymmddTime.Weekday(); weekday == time.Saturday || weekday == time.Sunday {
		return "", errors.New("六/日休市")
	}
	var filedir string = "downloads/"
	var filename string = "dailycorptpex_"
	var sub string = ".csv"
	url, _ := web.AppConfig.String("tpexcorpurl")
	url += yyyymmddTime.AddDate(-1911, 0, 0).Format("2006/01/02")[1:]

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
	} else if wr <= 512 {
		return "", errors.New("無內容")
	}
	fmt.Println("每日上櫃三大法人下載結束")
	return filedir + filename + yyyymmdd + sub, nil
}

func dailyCorpFilterTPEX(filename string) ([]common.CorporationDailyTrans, error) {
	fmt.Println("每日上櫃三大法人處理開始")
	corpdailys := []common.CorporationDailyTrans{}
	by, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	//字元編碼
	r, _, _ := transform.String(traditionalchinese.Big5.NewDecoder(), string(by))
	arr := strings.Split(r, "三大法人買賣超股數合計")

	//循環處理
	sc := bufio.NewScanner(strings.NewReader(strings.TrimSpace(arr[1])))
	for sc.Scan() {
		str := string(sc.Bytes())
		str = str[1 : len(str)-1]           //去除尾部符號以便分割
		strArr := strings.Split(str, `","`) //分割

		//數據轉換為模型
		corpdaily := common.CorporationDailyTrans{}
		corpdaily.Date, _ = strconv.ParseInt(filename[24:32], 10, 64)
		corpdaily.Code = strArr[0]
		corpdaily.Name = strArr[1]

		fiTemp, _ := strconv.Atoi(strings.Replace(strArr[4], ",", "", -1))
		corpdaily.ForeignInvestors = fiTemp / 1000

		fcTemp, _ := strconv.Atoi(strings.Replace(strArr[7], ",", "", -1))
		corpdaily.ForeignCorp = fcTemp / 1000

		itTemp, _ := strconv.Atoi(strings.Replace(strArr[13], ",", "", -1))
		corpdaily.InvestmentTrust = itTemp / 1000

		dsTemp, _ := strconv.Atoi(strings.Replace(strArr[16], ",", "", -1))
		corpdaily.DealerSelf = dsTemp / 1000

		dhTemp, _ := strconv.Atoi(strings.Replace(strArr[19], ",", "", -1))
		corpdaily.DealerHedge = dhTemp / 1000

		dTemp, _ := strconv.Atoi(strings.Replace(strArr[22], ",", "", -1))
		corpdaily.Dealer = dTemp / 1000

		tTemp, _ := strconv.Atoi(strings.Replace(strArr[23], ",", "", -1))
		corpdaily.Total = tTemp / 1000

		corpdaily.CDUnion = corpdaily.Code + "-" + fmt.Sprint(corpdaily.Date)
		corpdailys = append(corpdailys, corpdaily)
	}
	// 刪檔
	err = os.Remove(filename)
	if err != nil {
		return nil, err
	}
	fmt.Println("每日上櫃三大法人處理結束")
	return corpdailys, nil
}
