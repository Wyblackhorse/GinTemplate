/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package rule

import (
	"github.com/gin-gonic/gin"
	"github.com/wangyi/GinTemplate/model"
	"github.com/wangyi/GinTemplate/tools"
)

type LoginRequest struct {
	//Username string `json:"username" binding:"required"`
	Username   string `form:"username"  binding:"required,min=5,max=10"` //账户名字
	Password   string `form:"password"  binding:"required,min=5,max=10"` //密码
	GoogleCode string `form:"code"  binding:"omitempty,min=6,max=6"`     //谷歌验证码
}

type MenusJson struct {
	Top model.Menu
	Sec []model.Menu
}

type MenusJsonData struct {
	Data []MenusJson
}

//返回 -101
func ReturnErr101(c *gin.Context, err string) {
	tools.JsonWrite(c, -101, []string{}, err)
}

//返回  200
func ReturnSuccess(c *gin.Context, success string) {
	tools.JsonWrite(c, 200, []string{}, success)
}

//返回  200 带result 数据
func ReturnSuccessData(c *gin.Context, data interface{}, success string) {
	tools.JsonWrite(c, 200, data, success)
}

//第三方充值订单回调
type RecordOrderBack struct {
	Code   int    `json:"Code"`
	Msg    string `json:"Msg"`
	Result struct {
		Data string `json:"Data"`
	} `json:"Result"`
}

//第三方充值订单回调(订单参数)
type RecordOrderBackParameter struct {
	PlatformOrder string `json:"PlatformOrder"`
	RechargeAddress string `json:"RechargeAddress"`
	Username string `json:"Username"`
	AccountOrders float64 `json:"AccountOrders"`
	AccountPractical float64 `json:"AccountPractical"`
	RechargeType string `json:"RechargeType"`
	BackURL string `json:"BackUrl"`
}
