package tools

import (
	"bytes"
	cryptoRand "crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"github.com/gin-gonic/gin"
	eeor "github.com/wangyi/GinTemplate/error"
	"go.uber.org/zap"
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

var publicKey = []byte(`
-----BEGIN RSA Public Key-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAoEixcgAPmLpHLEDh3P8e
GpxolNbGJoxbrNQU1kaRCTMiu5qTaLJsTb6SVh7J4yztLOSdvIwbC2YeyVW8fatx
3eQ4RX6/txdtm07ov1bmC9n6/caOeRz2Pq2ZOse3uFuSjpQbF/2oAv3E6zWq5tdH
wG89ZNj+igs5lme4S6Uy2OE2MsqV/kwGMdBcdTOld8ki3MTsoEeBg9+IoqRD6gqi
l9sZdoHf0ItVE85Rw2Gp1rMfeTUMW7W3SvKItB33978/PgVmUvKLwY9+xsvWmILB
ZgkIMjUZ9/98LsyhdOvElcWFkYX2f84PSIJ8roAJhUpGJnw05i1jykDM6wa3sZ05
1wIDAQAB
-----END RSA Public Key-----
`)

//加密
func RsaEncrypt(origData []byte) ([]byte, error) {
	//解密pem格式的公钥
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, eeor.OtherError("public key error")
	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//加密
	return rsa.EncryptPKCS1v15(cryptoRand.Reader, pub, origData)
}

var privateKeyForEveryOne = []byte(`
-----BEGIN RSA Private Key-----
MIIEpAIBAAKCAQEAy3My0diBGfkjFzq5UVl0SeOLSg/Lcmvhl9hEuRr6B7O7KV7g
hYCQR2tJcfrZQ1ehqPVS1jskNkoXKXfcAzeEgQrcFLYMVwuHVh2mu8imUoKY7fke
arU6MtJHi2bpxownZLJzurbzbeeWiZWj05HzCZVPjfAhxkTdC+kuZBJfF0Fc6xrl
XgbDslsEyEyKIGku7G47ZRmtJjDiUk+Bec7f9uhTbSWWu4ZO57S4fuA5K1qXf0Tm
w6fEiM5DRkfYGsmO+2x6AmHGwVhFw37k/UEpur3bkajK9fk3s6xtJjmLet3y+g6j
cdpPMr+sZdMFxfGhIqu0xPy2mNNZmANJni3/PQIDAQABAoIBAA8LL6DQr4sqHuwi
zX00biLgjnYlgNevHnlJ5psBYaecJKTEfTmh7gk5565j7BjMrAmASmXI7b6N7/SD
BmO+gS/Bi9CEPZlaIuG9Q4zzI0lKmuBN4W/mgq0rW1r1eyfRSUBq6Z/O02U3EKyP
whNs4Vm+DqniLb0pbmbpESMZMKrZaqXSRMQHfwp7wx8W+3dPEQAs3pO1SwEpfL8S
bx/gBaRUi5iVIjG52G7daoCCGAPzyzhGtKNkc/Yff9dnI7hGwS2/zysrOychDSKJ
+/XGlQb9Zwp8F7bz84wm3BACTCCjzvGd5NxIZqXoeR8VZ9LWpJSlRIMIskExhlp/
hFfZtckCgYEA8Vc/m8frL879ioxSkbcaFYDuvS0FFbAZWh9z2XTJGgSWQlQ0Ms5k
U/668c7VXOlLrbUQrvH+aAp4uok+o+hlZNTbFJhRXs2igfIKpC4Ak8tYkbrsxmpJ
gP9/LScLlgTmrrFSo6ADJU0oPSpC95HwvwyIao9fpauhFYZtBf5s7o8CgYEA187D
YOuRbmHKbYLCXe37grEGerOz4WUpSWkaFmlwZ14T8i/U0ps/qRiiOOdfrEnjAT3L
JkZ1GL0fBLqN2KAkD7XjcOECcVYQMN/5qwEys9udmbgj+I0u5RSeE8mYB0RP06/q
Cg8NocxH7Z5STjd3JA1gZwl974uwcN1mkdduW3MCgYEAyUOXmlRowCAAtRA8s6Rd
Ll2tuznWKbYIDm54cHrCUt5MaNhMB6qzZJDkWk/BA5DTOfPsC9ln7l/9OqLGCG8A
T8xrP4ufIE6hHXk6gpySgq5sGGwolXeCAQARkRgkw2Em97yNTENfHDZyPkAGROwC
N3E+Oo+ClmjBF3BZb0w0j+UCgYA5Y46pc3uVMwQ14xP1DphXxOPINYmcYt572ytI
0nlFw8riGL4r04U2XoqlP0I9+tgXOGuRniL9lS1ugH3AIbX1R5VYKz4PDaf4l1c5
lnP5SGm8uy81pbXWzYjMEkwPgqcH0DwYuLATWtO16OhSTIWuXLBKNkf7L9aX7Qid
uABs6QKBgQCU6GhHoaIPAvVaMCm+sdT8MwH6ThH3a+kWqFj3AgkUBSv9UUhYQ+A6
Qey+AcRpKVP91z82noHucakJ/7EMTOnsnLK8mjrFak+L2GZWrPhk2bMZnF7vxGHu
FcVT4XWqDGY3GwdDTfX0WcJGtjyLJkD05ldO74jNYXVvIZvyN0abLw==
-----END RSA Private Key-----
`)

// RsaDecryptForEveryOne RsaDecrypt 解密
func RsaDecryptForEveryOne(ciphertext []byte) ([]byte, error) {
	//解密
	block, _ := pem.Decode(privateKeyForEveryOne)
	if block == nil {
		return nil, eeor.OtherError("private key error!")

	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err

	}
	// 解密
	return rsa.DecryptPKCS1v15(cryptoRand.Reader, priv, ciphertext)
}

//发送数据POST

func PostHttp(url string, data []byte) (string, error) {
	reader := bytes.NewReader(data)
	post, err := http.Post(url, "", reader)
	if err != nil {
		zap.L().Debug("PostHttp postUrl:" + url + "err:" + err.Error() + "data:" + string(data))
		return "", err
	}
	respBytes, err1 := ioutil.ReadAll(post.Body)
	if err1 != nil {
		zap.L().Debug("PostHttp postUrl:" + url + "err:" + err.Error() + "data:" + string(data))
		return "", err
	}
	zap.L().Debug("PostHttp postUrl:" + url + "OK:" + string(respBytes) + "data:" + string(data))
	return string(respBytes), nil
}

// BackUrlToPay 第三方支付回调方法
func BackUrlToPay(backUrl string, bytesData string) (bool, error) {
	type TT struct {
		Code   int
		Msg    string
		Result struct {
			Data string
		}
	}
	var tt TT
	tt.Code = 200
	tt.Msg = "success"
	tt.Result.Data = bytesData
	data, err := json.Marshal(tt)
	fmt.Println(string(data))
	if err != nil {
		return false, err
	}
	reader := bytes.NewReader(data)
	post, err := http.Post(backUrl, "", reader)

	if err != nil {
		zap.L().Debug("回调地址:" + backUrl + "错误:" + err.Error())
		return false, err
	}
	respBytes, err1 := ioutil.ReadAll(post.Body)
	if err1 != nil {
		zap.L().Debug("回调地址:" + backUrl + "错误:" + err1.Error())
		return false, err1
	}
	zap.L().Debug("回调地址:" + backUrl + " 返回结果:" + string(respBytes))
	return true, nil
}
