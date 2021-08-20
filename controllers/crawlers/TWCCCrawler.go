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

	"github.com/beego/beego/v2/server/web"
	common "github.com/bensan0/stockinfo/models/common"
)

//集保所持股分布

func TWCCDCrawl() error {
	filename, err := DownloadTWCCDistribution()
	if len(filename) == 0 || err != nil {
		return err
	}
	err = TWCCDFilter(filename)
	if err != nil {
		return err
	}
	return nil
}

//集保所下載持股分布
func DownloadTWCCDistribution() (string, error) {
	fmt.Println("開始下載股權分佈")
	url, _ := web.AppConfig.String("distributionurl")
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	//create
	out, err := os.Create("downloads/TWCCDistribution.csv")
	if err != nil {
		return "", err
	}
	defer out.Close()

	//copy
	wr, err := io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	} else if wr == 0 {
		err := errors.New("無內容")
		return "", err
	}
	fmt.Println("股權分佈下載結束")
	return "downloads/TWCCDistribution.csv", nil
}

func TWCCDFilter(filename string) error {
	var ownerships = []common.OwnershipDistribution{}
	var ownership common.OwnershipDistribution
	fmt.Println("開始解析股權分佈")
	by, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	strArr := strings.Split(string(by), "占集保庫存數比例%")
	sc := bufio.NewScanner(strings.NewReader(strings.TrimSpace(strArr[1])))

	for sc.Scan() {
		str := string(sc.Bytes())
		strArr := strings.Split(str, ",")
		rate, _ := strconv.ParseInt(strArr[2], 10, 64)
		date, _ := strconv.ParseInt(strArr[0], 10, 64)
		code := strArr[1]
		popu := strArr[3]
		vol, _ := strconv.ParseInt(strArr[4], 10, 64)
		vol = vol / 1000
		percentage := strArr[5]

		switch rate {
		case 1:
			ownership = common.OwnershipDistribution{}
			ownership.Rate1 = popu + "/" + fmt.Sprint(vol) + "/" + percentage
			ownership.Code = code
			ownership.Date = date
			ownership.CDUnion = code + "-" + fmt.Sprint(date)
		case 2:
			ownership.Rate2 = popu + "/" + fmt.Sprint(vol) + "/" + percentage
		case 3:
			ownership.Rate3 = popu + "/" + fmt.Sprint(vol) + "/" + percentage
		case 4:
			ownership.Rate4 = popu + "/" + fmt.Sprint(vol) + "/" + percentage
		case 5:
			ownership.Rate5 = popu + "/" + fmt.Sprint(vol) + "/" + percentage
		case 6:
			ownership.Rate6 = popu + "/" + fmt.Sprint(vol) + "/" + percentage
		case 7:
			ownership.Rate7 = popu + "/" + fmt.Sprint(vol) + "/" + percentage
		case 8:
			ownership.Rate8 = popu + "/" + fmt.Sprint(vol) + "/" + percentage
		case 9:
			ownership.Rate9 = popu + "/" + fmt.Sprint(vol) + "/" + percentage
		case 10:
			ownership.Rate10 = popu + "/" + fmt.Sprint(vol) + "/" + percentage
		case 11:
			ownership.Rate11 = popu + "/" + fmt.Sprint(vol) + "/" + percentage
		case 12:
			ownership.Rate12 = popu + "/" + fmt.Sprint(vol) + "/" + percentage
		case 13:
			ownership.Rate13 = popu + "/" + fmt.Sprint(vol) + "/" + percentage
		case 14:
			ownership.Rate14 = popu + "/" + fmt.Sprint(vol) + "/" + percentage
		case 15:
			ownership.Rate15 = popu + "/" + fmt.Sprint(vol) + "/" + percentage
		case 16:
			continue
		case 17:
			ownership.Total = popu + "/" + fmt.Sprint(vol) + "/" + percentage
			ownerships = append(ownerships, ownership)
		}

	}
	//刪檔
	err = os.Remove(filename)
	if err != nil {
		return err
	}
	fmt.Println("解析股權分佈結束")
	//持久化
	result := common.DB.Create(ownerships)
	if result.Error != nil {
		fmt.Println("存入DB失敗")
		return result.Error
	}
	fmt.Println("存入DB結束")
	return nil
}
