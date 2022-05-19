package tools

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"testing"
)

func TestGetRootPath(t *testing.T) {
	response, err := http.Get("http://www.haxi06666.com")
	if err != nil {
		fmt.Println("---")
		fmt.Println(err.Error())
		return
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	fmt.Println(string(body))
}
