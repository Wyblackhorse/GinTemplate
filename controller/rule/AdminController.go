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
	"github.com/wangyi/GinTemplate/dao/mysql"
	"github.com/wangyi/GinTemplate/dao/redis"
	"github.com/wangyi/GinTemplate/model"
	"github.com/wangyi/GinTemplate/tools"
	"time"
)

//登录
func Login(c *gin.Context) {
	var admin LoginRequest
	if err := c.ShouldBind(&admin); err != nil {
		tools.JsonWrite(c, -101, []string{}, "err:"+err.Error())
		return
	}
	adminTwo := model.AdminModel{}
	err := mysql.DB.Where("username=?", admin.Username).Where("password=?", admin.Password).First(&adminTwo).Error
	if err != nil {
		tools.JsonWrite(c, -101, []string{}, "账户或者密码错误!")
		return
	}
	//查看账号状态
	if adminTwo.Status == 2 {
		tools.JsonWrite(c, -101, []string{}, "该账号已经停用")
		return
	}
	//if adminTwo.GoogleCode == "" {
	//	secret, _, qrCodeUrl := tools.InitAuth(adminTwo.Username)
	//	err := mysql.DB.Model(&model.AdminModel{}).Where("id=?", adminTwo.ID).Update(model.AdminModel{GoogleCode: secret}).Error
	//	if err != nil {
	//		tools.JsonWrite(c, -101, []string{}, "err:"+err.Error())
	//		return
	//	}
	//	tools.JsonWrite(c, -102, map[string]string{"codeUrl": qrCodeUrl}, "请先绑定谷歌账号")
	//	return
	//}
	////校验谷歌验证
	//verifyCode, _ := tools.NewGoogleAuth().VerifyCode(adminTwo.GoogleCode, admin.GoogleCode)
	//if !verifyCode {
	//	tools.JsonWrite(c, -101, map[string]string{}, "谷歌验证失败")
	//	return
	//}
	//登录成功
	redis.Rdb.Set("Admin_Login_Token_"+adminTwo.Token, adminTwo.Username, time.Second*3600*24)
	//登录日志
	log := model.Loggers{Content: "登陆成功", WriteId: int(adminTwo.ID)}
	log.AdminAddErrorLoggers(mysql.DB)
	tools.JsonWrite(c, 200, adminTwo.Token, "登录成功")
	return
}

