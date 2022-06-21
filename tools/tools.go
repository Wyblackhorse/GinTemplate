package tools

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"image/gif"
	"image/jpeg"
	"image/png"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
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

//生成随机字符串
func RandStringRunes(n int) string {
	var letterRunes = []rune("0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}

	return string(b)
}

/**
数组是否存在摸一个元素
*/
func InArray(arr []string, param string) bool {
	for _, v := range arr {
		if param == v {
			return true
		}
	}
	return false
}

//json返回
func JsonWrite(context *gin.Context, status int, result interface{}, msg string) {
	context.JSON(http.StatusOK, gin.H{
		"code":   status,
		"result": result,
		"msg":    msg,
	})
}

//返回当前是第几周

func ReturnTheWeek() int {
	datetime := time.Now().Format("20060102")
	timeLayout := "20060102"
	loc, _ := time.LoadLocation("Local")
	tmp, _ := time.ParseInLocation(timeLayout, datetime, loc)
	_, intWeek := tmp.ISOWeek()
	return intWeek
}

func ReturnTheMonth() int {
	_, m, _ := time.Now().Date()
	return int(m)
}

//获取up主的视频链接   youTuBe   返回数组  url
func FormYouTuBeUrl(url string) []string {
	res, err := http.Get(url)

	if err != nil {
		fmt.Println(err.Error())
		return []string{}
	}
	defer res.Body.Close()
	req, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error())
		return []string{}

	}
	//watch\?v=\S{11}
	reg1 := regexp.MustCompile(`watch\?v=\S{11}`)
	if reg1 == nil { //解释失败，返回nil
		fmt.Println("regexp err")
		return []string{}
	}
	//根据规则提取关键信息
	result1 := reg1.FindAllStringSubmatch(string(req), -1)
	if len(result1) < 1 {
		return []string{}
	}
	var p []string
	for _, v := range result1 {
		p = append(p, "https://www.youtube.com/"+v[0])
	}

	fmt.Println(1312312)
	return p
}

//检查木马文件
func CheckImageFile(path, style string) (string, error) {
	f, err := os.Open(path)

	fmt.Println(f)
	if err != nil {
		_ = fmt.Errorf("打开文件失败 %s", err.Error())
	}
	switch strings.ToUpper(style) {
	case "JPG", "JPEG":
		_, err = jpeg.Decode(f)
	case "PNG":
		_, err = png.Decode(f)
		fmt.Println(err)
	case "GIF":
		_, err = gif.Decode(f)
	}
	if err != nil {
		_ = fmt.Errorf("校验文件类型失败 %s", err.Error())
		return "", err
	}
	return "", nil
}
