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


  fmt.Println(	time.Now().AddDate(0,-1,0).Format("2006-01-02"))


  fmt.Println(time.Now().Date())
	_,m,_:=time.Now().Date()
	fmt.Println(int(m))

	fmt.Println(time.Now().Weekday())
	datetime := time.Now().Format("20060102")
	timeLayout := "20060102"
	loc, _ := time.LoadLocation("Local")
	tmp, _ := time.ParseInLocation(timeLayout, datetime, loc)
	_, intWeek := tmp.ISOWeek()
	fmt.Println(intWeek)
}



func TestTwo(t *testing.T){

fmt.Println(time.Now().Format("2006-01-02"))








}