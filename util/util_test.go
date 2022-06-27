package util

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func TestName(t *testing.T) {

	type Create struct {
		PlatformOrder    string
		RechargeAddress  string
		Username         string
		AccountOrders    float64 //订单充值金额
		AccountPractical float64 //  实际充值的金额
		RechargeType     string
		BackUrl          string
	}

	p := Create{
		PlatformOrder:    "2012553332254545254224252455",
		Username:         "wing",
		RechargeType:     "USDT",
		AccountOrders:    200.00,
		AccountPractical: 200.00,
		BackUrl:          "https://123.com",
	}

	//c:=Stu{Name: "西欧奥课啊",Age: 10}
	data, err := json.Marshal(p)
	fmt.Println(string(data))
	data, _ = RsaEncryptForEveryOne(data)
	fmt.Println(err)
	fmt.Println(data)
	fmt.Println(base64.StdEncoding.EncodeToString(data))




	pay, err := BackUrlToPay("http://127.0.0.1:7777/v2/backUrl", base64.StdEncoding.EncodeToString(data))
	if err != nil {
		return
	}
	fmt.Println(pay)
}

func TestBackUrlToPay(t *testing.T) {
	type Create struct {
		PlatformOrder    string
		RechargeAddress  string
		Username         string
		AccountOrders    float64 //订单充值金额
		AccountPractical float64 //  实际充值的金额
		RechargeType     string
		BackUrl          string
	}

	var tt Create
	tt.PlatformOrder = "1111111111111111111111111111"
	tt.RechargeAddress = "TW2HWaLWy9pwiRN4yLju6YKW3aQ6Fw8888"
	tt.Username = "wine"
	tt.AccountOrders = 100.00
	tt.AccountPractical = 10.00
	tt.RechargeType = "USDT"
	data, err := json.Marshal(tt)

	if err != nil {
		fmt.Println(err.Error())
	}
	data, _ = RsaEncryptForEveryOne(data)

	BackUrlToPay("http://8.136.97.179:7777/v2/backUrl", base64.StdEncoding.EncodeToString(data))
}

func TestBackUrl11ToPay(f *testing.T) {
	type Create struct {
		Username string
	}
	var tt Create
	tt.Username = "jack"
	data, err := json.Marshal(tt)
	if err != nil {
		fmt.Println(err.Error())
	}

	data, err = RsaEncrypt(data)
	fmt.Println(base64.StdEncoding.EncodeToString(data))
}

func TestRsaEncrypt(t *testing.T) {

	fmt.Println(time.Now().Unix())
}
