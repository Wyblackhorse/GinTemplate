/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package tools

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"testing"
	"time"
)

func TestRandStringRunes(t *testing.T) {

	fmt.Println("-----------------开启二次认证----------------------")
	user := "testxxx@qq.com"
	secret, qrCodeUrl, a := InitAuth(user)

	fmt.Println("-----------------信息校验----------------------")

	fmt.Println(secret, qrCodeUrl, a)

	// secret最好持久化保存在
	//// 验证,动态码(从谷歌验证器获取或者freeotp获取)
	//bool, err := NewGoogleAuth().VerifyCode("QNSCBEQ4M6MAOIWB3H2NZYMN3JUD3SVT", "899397")
	//if bool {
	//	fmt.Println("√")
	//} else {
	//	fmt.Println("X", err)
	//}

}

func TestOne(t *testing.T) {

	fmt.Println(time.Now().AddDate(0, -1, 0).Format("2006-01-02"))

	fmt.Println(time.Now().Date())
	_, m, _ := time.Now().Date()
	fmt.Println(int(m))

	fmt.Println(time.Now().Weekday())
	datetime := time.Now().Format("20060102")
	timeLayout := "20060102"
	loc, _ := time.LoadLocation("Local")
	tmp, _ := time.ParseInLocation(timeLayout, datetime, loc)
	_, intWeek := tmp.ISOWeek()
	fmt.Println(intWeek)
}

func TestTwo(t *testing.T) {

	type ReturnData struct {
		TaskNum        int     `json:"task_num"`        //任务数量
		TaskEarnings   float64 `json:"task_earnings"`   //任务收益
		JuniorTaskNum  int     `json:"junior_task_num"` //下级任务数量
		JuniorEarnings float64 `json:"junior_earnings"` //下级任务收益
		Date           string  `json:"date"`
	}

	//c := make([]ReturnData, 0)
	//c1 := ReturnData{TaskNum: 10, TaskEarnings: 100, Date: "6-16"}
	//c2 := ReturnData{TaskNum: 9, TaskEarnings: 10, Date: "6-12"}
	//c3 := ReturnData{TaskNum: 8, TaskEarnings: 190, Date: "6-13"}
	//
	//
	//
	//b1 := ReturnData{JuniorTaskNum: 10, JuniorEarnings: 100, Date: "6-16"}
	//b2 := ReturnData{JuniorTaskNum: 9, JuniorEarnings: 10, Date: "6-12"}
	//b3 := ReturnData{JuniorTaskNum: 8, JuniorEarnings: 190, Date: "6-10"}
	//b4 := ReturnData{JuniorTaskNum: 8, JuniorEarnings: 200, Date: "6-17"}

}

func TestFormYouTuBeUrl(t *testing.T) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://www.tiktok.com/@cicichania96", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("authority", "www.tiktok.com")
	req.Header.Set("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	req.Header.Set("accept-language", "zh-CN,zh;q=0.9")
	req.Header.Set("cache-control", "max-age=0")
	req.Header.Set("cookie", "tt_csrf_token=Y2bcQCyE-KW576loLSvuBQ9lXC2B0ZUBVcPQ; _abck=4D640AE24F04EC4CDF2C35618165BFA1~-1~YAAQT8cGF/nuwriBAQAAYJ2tGgh+FU5gxp39lm+s6WpjcSJrOL1FOOaYlzHnjUgwlRWXqWeskAHaudeLvcE1QIjbgAP/nHvZAlHeeHESwRyRmSoogWU4fLH3gGX59wl6YF4ewzf25MJegY9bl6Le4NP3HFDwsa8A4L/ONxHTux+mFqxxK6edOsdMtxRcjqjYiZ1d1EERDeg9VQgsMEjpQwaNNnZ7vVBnoYhoOR2TQ7K1WZCezzdJIm5v+iEvHxu1PAFgKHJE5MPgQ2i0iR2pTrf4Ta9m9u0tbWnPMS3boKa0fA55xEUcVd9jYnXqUOkY1ClCDqsWSOrYFktJsD247IKYrIyVoh3PReajynUvbPVcv9b9zmu07hN5LdI=~-1~-1~-1; bm_sz=3939BD3B4492D06D163070D7CA20E6F2~YAAQT8cGF/ruwriBAQAAYJ2tGhDqtj2i9+ASn9HnTxn5w7wkTbiSPVgoh0DKq9LNQEHg/GU7Z8/16seScGxvahjvYGxHBucUNmsL+tz8mMgGY5rN39sBgNOI/otwfyoTy0KTqRHdw5Jf41c4L8u5ba8Di+VPFeNgcCNwf1qYgW02im3OCMz3xnZUSGsEpS2aMpnvx25+2aniCHaTor2J7aBMtDCS9wTEvnLoQfKy6ZSh4fY8DsqMYPKRaMjWtHvoq6ouq08alvMCp3TN3P6hWI6rNv4YQEE3Hum21Gn3VMpPB9U=~3159618~3556405; ak_bmsc=1376967A00391DDBAA013A80BE7970C8~000000000000000000000000000000~YAAQBLdNaCkwa7WBAQAAH+qtGhDNYfA3JNPBM6c65o0OK9NGJq4NHrEF/u7j+rVN6fseLHHFwnK9P1tMrBabmTvCIFk3M5QlyrUmkJs/4r6P9Or7kKniRka1/gRo+8752ljjdH6Yvbj+CXYV6FxgCdLFL4WP6uBikg9IgO2YbsUZ7T14AztdoSM4o3qzLpzCSwNQoFYrYngUzz9iPRFqQKqV6dyMJw0zaWhfvIYXBJyafEZJFps6HQ/kRONajB0d0hcCgRC2c0jF94cDCUJU3IpsWysSW5BKF65xCsLG+9fRQhl5f6IAM3ZW7OJGJIpXpvdFFVEJ7qFQ8czyVOgu5zsUUvVpODEh2a6jFK/v5sZkpcPimVI9vdtIPL4twWbQTDVLoZph6pekMA==; __tea_cache_tokens_1988={%22_type_%22:%22default%22%2C%22user_unique_id%22:%227122365539294430766%22%2C%22timestamp%22:1658305104597}; s_v_web_id=verify_4d69ed99aa27361f8aa1f88afe094ca9; ttwid=1%7CiR6hxOtUzuoq6bjg53FiUrZszCXaJqovM2DAVLSw3sM%7C1658305192%7Cda3f6c3dede7792b51122a7986d4f8bd50da2abc52a4284bef7fd5bdca0660ea; bm_sv=F81FF4D4B505C11AB564CC0719CF6FE5~YAAQGi0+F8ODQBSCAQAAVy6xGhAk1zK/OcoZ7+tS77Guk2cVz+cos9DUu7IQ/qEEbgrQe5NEk+kFtm6/HRqkFZ/pXGHSgA4+ucUI+rbrYU+oxU81zw26f0axLn1ZCJBYZx3p3pkLDatePpv5FNa3z5FA71j3v2XIb3y9JR1qG6Hn4SIo2m5fGHUltEDea4rYnO7b+pRQKb57kjE2mcdwZHfRQDhcq6nQPHXAujJUV5tXilThHzG75KHoo2wtwzNr~1; msToken=Mb0BCIhNi70FjCaNdLfku3rtnwJ1RHCyKygYENg6ASW1RjNEnGUm2r5N2MLx1eBQAHz3cZ0e8yGCz30rMab7DzbQ0wiEfDMcIKB1fHj8x1kL56pC1fcN; msToken=Mb0BCIhNi70FjCaNdLfku3rtnwJ1RHCyKygYENg6ASW1RjNEnGUm2r5N2MLx1eBQAHz3cZ0e8yGCz30rMab7DzbQ0wiEfDMcIKB1fHj8x1kL56pC1fcN")
	req.Header.Set("sec-ch-ua", `".Not/A)Brand";v="99", "Google Chrome";v="103", "Chromium";v="103"`)
	req.Header.Set("sec-ch-ua-mobile", "?0")
	req.Header.Set("sec-ch-ua-platform", `"macOS"`)
	req.Header.Set("sec-fetch-dest", "document")
	req.Header.Set("sec-fetch-mode", "navigate")
	req.Header.Set("sec-fetch-site", "none")
	req.Header.Set("sec-fetch-user", "?1")
	req.Header.Set("upgrade-insecure-requests", "1")
	req.Header.Set("user-agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/103.0.0.0 Safari/537.36")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	//解析正则表达式，如果成功返回解释器
	reg1 := regexp.MustCompile(`{"user-post":{"list":.*"],"browserList"`)
	if reg1 == nil {
		fmt.Println("regexp err")
		return
	}
	//根据规则提取关键信息
	result1 := reg1.FindAllStringSubmatch(string(bodyText), -1)
	if len(result1) > 0 {
		if len(result1[0]) > 0 {
			reg1 = regexp.MustCompile(`\d+`)
			result1 = reg1.FindAllStringSubmatch(result1[0][0], -1)
			for _, i2 := range result1 {
				fmt.Println(i2[0])
			}
		}
	}

}
