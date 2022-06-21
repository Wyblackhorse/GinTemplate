package util

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

//私钥
var privateKey = []byte(`
-----BEGIN RSA Private Key-----
MIIEpAIBAAKCAQEAoEixcgAPmLpHLEDh3P8eGpxolNbGJoxbrNQU1kaRCTMiu5qT
aLJsTb6SVh7J4yztLOSdvIwbC2YeyVW8fatx3eQ4RX6/txdtm07ov1bmC9n6/caO
eRz2Pq2ZOse3uFuSjpQbF/2oAv3E6zWq5tdHwG89ZNj+igs5lme4S6Uy2OE2MsqV
/kwGMdBcdTOld8ki3MTsoEeBg9+IoqRD6gqil9sZdoHf0ItVE85Rw2Gp1rMfeTUM
W7W3SvKItB33978/PgVmUvKLwY9+xsvWmILBZgkIMjUZ9/98LsyhdOvElcWFkYX2
f84PSIJ8roAJhUpGJnw05i1jykDM6wa3sZ051wIDAQABAoIBAEOVo0zAfdMWaMBe
M16uLBFodiicTz0d3eIIIyke5DUO3MxiJ1n0MfquKqIppVMnNHLDi0dfhF/QFZcr
kakjy4WWn6ueAFTBijP8l+NmTuYcadrsCFNFRQe0x3GjmPIhmrCH/syk5l0siAXt
wTxI9jZMMYT+goXOqd+jqHAtHj2YPh1wWSXLfyX+M3NwOXQ9DlDCwZ8SJjxYl8KI
Q5CJoTjY36Qi0y7XlXUWPDJ1vBH5E6fkNqWcdAKDvxmFDscegbUo3gZ1md6PCV7T
0cpUXoYYEvvrWAwQae5gozb3/rIRXEC6z8iAY/XzwxHB/RnDivAMYufGhBhFtqGr
l3oFCcECgYEAzOLCc6rxuobaLR13Dshb3EvKn7O/KCGKJd/Izpgt33VcS1L4dPIH
L3HDZw8KuvysjXnE7nCC862NfAGqNDMz2lbLgkIdRRxDx64KTUeZCaoMdz8rexKT
VKDGLp5jJI9FU2EnbNG+irLPkbCcXEKEfNwwVNcGAIWDIiJTpe2nJacCgYEAyEVk
eMHF5ZozLxYwHs4ybQgtIzNiq3oF5m2z6zlSAk967an4jdpb3VS1xOtsdSBF2LW3
dnvV3mWu/35U7nkRh0tKPXNZc3HV3CKezzOxZ0zoAj9elwb9oriPWuKVisq18Ujt
a5QVOGBXUzc0c1ONjDBh0mvRNzhr7qOMig1dMFECgYBsimjXXCVJWq10nxp2o2A+
2YwThObs/K+yFtbL08Thj8wAP4lOcvWphcwt6cMWgktre6n/Y22MaFH+8ubXVpTO
w5J2hE37Udj6jNH6VMbXXtXRyo5fWdzhRXcYNWJyeNASNvLq7EbUNZxPI1ACdF65
wvB70ZnlZtWsnKDR04/sGwKBgQCn1U+Py4QXGJTQXx3Qkyi7KuD44PVNkyMiqsje
1dieSxFP3uOHrXjTEUyLTGhF99fQ9uhbCQiAKmLvhmWSvC8uXLBIs0RBdSKuKu0/
46hGU7MTPxv8IUWpelXY6o48FAlJvb4KK71k04gbGuZ/x4OV+m3gM67PQh9hi/oZ
L33rIQKBgQC7/nLHBWcrs2VPGivKUc9ZeHbFLdsPM6zGaZoarCaAFEYZn3DF9r5s
Gkn6NoLph+OO4qPla2YOwfZuT8MokQ7eHHNxZPVGMrbmdPPofzzoYNljar+4r7pk
TuIT9b1m+Ao/IA19g76Mr+uzQnJYZN2kxewqkI+x05yARTXPoSIRHg==
-----END RSA Private Key-----
`)

// 公钥: 根据私钥生成
//openssl rsa -in rsa_private_key.pem -pubout -out rsa_public_key.pem
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

var publicKeyForEveryOne = []byte(`
-----BEGIN RSA Public Key-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAy3My0diBGfkjFzq5UVl0
SeOLSg/Lcmvhl9hEuRr6B7O7KV7ghYCQR2tJcfrZQ1ehqPVS1jskNkoXKXfcAzeE
gQrcFLYMVwuHVh2mu8imUoKY7fkearU6MtJHi2bpxownZLJzurbzbeeWiZWj05Hz
CZVPjfAhxkTdC+kuZBJfF0Fc6xrlXgbDslsEyEyKIGku7G47ZRmtJjDiUk+Bec7f
9uhTbSWWu4ZO57S4fuA5K1qXf0Tmw6fEiM5DRkfYGsmO+2x6AmHGwVhFw37k/UEp
ur3bkajK9fk3s6xtJjmLet3y+g6jcdpPMr+sZdMFxfGhIqu0xPy2mNNZmANJni3/
PQIDAQAB
-----END RSA Public Key-----
`)

// RsaEncrypt 加密
func RsaEncrypt(origData []byte) ([]byte, error) {
	//解密pem格式的公钥
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")

	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err

	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//加密
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// RsaDecrypt 解密
func RsaDecrypt(ciphertext []byte) ([]byte, error) {
	//解密
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")

	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err

	}
	// 解密
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

func RsaEncryptForEveryOne(origData []byte) ([]byte, error) {
	//解密pem格式的公钥
	block, _ := pem.Decode(publicKeyForEveryOne)
	if block == nil {
		return nil, errors.New("public key error")

	}
	// 解析公钥
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err

	}
	// 类型断言
	pub := pubInterface.(*rsa.PublicKey)
	//加密
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// RsaDecryptForEveryOne RsaDecrypt 解密
func RsaDecryptForEveryOne(ciphertext []byte) ([]byte, error) {
	//解密
	block, _ := pem.Decode(privateKeyForEveryOne)
	if block == nil {
		return nil, errors.New("private key error!")

	}
	//解析PKCS1格式的私钥
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err

	}
	// 解密
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
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

	//type T struct {
	//	Code   int    `json:"Code"`
	//	Msg    string `json:"Msg"`
	//	Result struct {
	//		Data string `json:"Data"`
	//	} `json:"Result"`
	//}
	//var jsonData T
	//err3 := json.Unmarshal(respBytes, &jsonData)
	//if err3 != nil {
	//	zap.L().Debug("回调地址:" + backUrl + "错误:" + err3.Error())
	//	return false, err3
	//}
	//if jsonData.Code != 200 {
	//	zap.L().Debug("回调地址:" + backUrl + "错误:" + jsonData.Msg)
	//	return false, eeor.OtherError("code不等于200 " + jsonData.Msg)
	//}
	//
	//zap.L().Debug("回调地址:" + backUrl + "成功")
	return true, nil
}
