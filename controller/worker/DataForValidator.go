/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package worker

import (
	"github.com/gin-gonic/gin"
	"github.com/wangyi/GinTemplate/tools"
)

//注册验证
type RegisterRequest struct {
	//	Username   string `form:"username"  binding:"required,min=5,max=10"` //账户名字
	Phone            string `form:"phone"  binding:"required,max=20"`
	Password         string `form:"password"  binding:"required,max=20"`
	PasswordTwo      string `form:"password_two"  binding:"required,max=20"`
	VerificationCode string `form:"verification_code"  binding:"required,min=4,max=4"` //验证码
	InvitationCode   string `form:"invitation_code"  binding:"omitempty,min=6,max=6"`  //邀请码
	VerificationId   string `form:"verification_id"  binding:"required"`
}

//登录验证
type LoginRequest struct {
	Phone    string `form:"phone"  binding:"required,max=20"`
	Password string `form:"password"  binding:"required,max=20"`
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
