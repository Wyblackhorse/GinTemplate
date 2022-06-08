package util

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	type Create struct {
		PlatformOrder   string
		RechargeAddress string
		Username        string
		AccountOrders   float64
		RechargeType    string
		BackUrl         string
	}
	p := Create{
		PlatformOrder:   "2012553332254545254224252455",
		RechargeAddress: "TW2HWaLWy9pwiRN4yLju6YKW3aQ6Fw8888",
		Username:        "wing",
		RechargeType:    "USDT",
		AccountOrders:   200.00,
		BackUrl:         "https://123.com",
	}

	//c:=Stu{Name: "西欧奥课啊",Age: 10}
	data, err := json.Marshal(p)
	data, _ = RsaEncryptForEveryOne(data)
	fmt.Println(err)
	fmt.Println(data)
	fmt.Println(base64.StdEncoding.EncodeToString(data))
	pay, err := BackUrlToPay("http://8.136.97.179:7777/v2/backUrl", base64.StdEncoding.EncodeToString(data))
	if err != nil {
		return
	}
	fmt.Println(pay)
}

