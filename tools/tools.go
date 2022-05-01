package tools

import (
	"github.com/gin-gonic/gin"
	"math/big"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"github.com/shopspring/decimal"
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

//判断文件文件夹是否存在(字节0也算不存在)
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
