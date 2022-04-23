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
)

func TestRandStringRunes(t *testing.T) {

	fmt.Println("-----------------开启二次认证----------------------")
	user := "testxxx@qq.com"
	secret, qrCodeUrl, a := InitAuth(user)



	fmt.Println("-----------------信息校验----------------------")

	fmt.Println(secret, qrCodeUrl,a)


	// secret最好持久化保存在
	//// 验证,动态码(从谷歌验证器获取或者freeotp获取)
	//bool, err := NewGoogleAuth().VerifyCode("QNSCBEQ4M6MAOIWB3H2NZYMN3JUD3SVT", "899397")
	//if bool {
	//	fmt.Println("√")
	//} else {
	//	fmt.Println("X", err)
	//}

}

