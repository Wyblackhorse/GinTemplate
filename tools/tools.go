package tools

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
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

func JsonWrite(c *gin.Context, code int, result interface{}, msg string) {
	c.JSON(http.StatusOK, gin.H{
		"code":   code,
		"result": result,
		"msg":    msg,
	})

}

func ReturnFail101(c *gin.Context, msg string) {
	JsonWrite(c, -101, map[string]interface{}{}, msg)
}
func ReturnSuccess200(c *gin.Context, msg string) {
	JsonWrite(c, 200, map[string]interface{}{}, msg)
}
func ReturnData200(c *gin.Context, result interface{}, msg string) {
	JsonWrite(c, 200, result, msg)
}
