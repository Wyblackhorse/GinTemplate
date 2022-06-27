package tools

import (
	"fmt"
	"github.com/agclqq/goencryption"
	"testing"
)

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

func TestApiSign(t *testing.T) {
	GET()
}