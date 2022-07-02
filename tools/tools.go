package tools

import (
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"math/big"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// GetRunPath2 获取程序执行目录
func GetRunPath2() string {
	file, _ := exec.LookPath(os.Args[0])
	path, _ := filepath.Abs(file)
	index := strings.LastIndex(path, string(os.PathSeparator))
	ret := path[:index]
	return ret
}

// IsFileNotExist 判断文件文件夹不存在
func IsFileNotExist(path string) (bool, error) {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return true, nil
	}
	return false, err
}

// IsFileExist 判断文件文件夹是否存在(字节0也算不存在)
func IsFileExist(path string) (bool, error) {
	fileInfo, err := os.Stat(path)

	if os.IsNotExist(err) {
		return false, nil
	}
	//我这里判断了如果是0也算不存在
	if fileInfo.Size() == 0 {
		return false, nil
	}
	if err == nil {
		return true, nil
	}
	return false, err
}

// GetRootPath 获取程序根目录
func GetRootPath() string {
	rootPath, _ := os.Getwd()
	if notExist, _ := IsFileNotExist(rootPath); notExist {
		rootPath = GetRunPath2()
		if notExist, _ := IsFileNotExist(rootPath); notExist {
			rootPath = "."
		}
	}
	return rootPath
}

func ReturnError101(context *gin.Context, msg string) {
	context.JSON(http.StatusOK, gin.H{
		"code":   -101,
		"result": map[string]string{},
		"msg":    msg,
	})
}
func ReturnError200(context *gin.Context, msg string) {
	context.JSON(http.StatusOK, gin.H{
		"code":   200,
		"result": map[string]string{},
		"msg":    msg,
	})
}
func ReturnError200Data(context *gin.Context, result interface{},msg string) {
	context.JSON(http.StatusOK, gin.H{
		"code":   200,
		"result": result,
		"msg":    msg,
	})
}

func ToDecimal(ivalue interface{}, decimals int) decimal.Decimal {
	value := new(big.Int)
	switch v := ivalue.(type) {
	case string:
		value.SetString(v, 10)
	case *big.Int:
		value = v
	}

	mul := decimal.NewFromFloat(float64(10)).Pow(decimal.NewFromFloat(float64(decimals)))
	num, _ := decimal.NewFromString(value.String())
	result := num.Div(mul)

	return result
}

func RandString(length int) []byte {
	var strByte = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890")
	var strByteLen = len(strByte)
	bytes := make([]byte, length)
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		bytes[i] = strByte[r.Intn(strByteLen)]
	}

	return bytes
}

//var apiKey = "6e347830f78d54482097d177f7278fb8"

func ApiSign(data ...[]byte) string {
	h256 := sha256.New()
	for _, b := range data {
		h256.Write(b)
	}
	hash := h256.Sum(nil)
	hashHex := hex.EncodeToString(hash)
	return hashHex
}

func HttpRequest(url string, params map[string]interface{},apiKey string) ([]byte, error) {
	data, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(data))
	req := make(map[string]interface{})
	b64Data := base64.StdEncoding.EncodeToString(data)
	req["data"] = b64Data
	req["sign"] = ApiSign([]byte(apiKey), []byte(b64Data), []byte(apiKey))
	reqData, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}
	payload := strings.NewReader(string(reqData))
	httpReq, err := http.NewRequest("POST", url, payload)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Add("Accept", "application/json")
	httpReq.Header.Add("Content-Type", "application/json")
	res, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return nil, errors.New("http status code:" + strconv.FormatInt(int64(res.StatusCode), 10))
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

