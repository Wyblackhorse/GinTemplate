/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package main

import (
	"fmt"
	"github.com/agclqq/goencryption"
	"github.com/wangyi/GinTemplate/cmd"
)

func main() {
	cmd.Execute()


}

func GET() {
	priv, err := goencryption.GenPrvKey(2048)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", priv)
	pub, err := goencryption.GenPubKeyFromPrvKey(priv)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", pub)
}
