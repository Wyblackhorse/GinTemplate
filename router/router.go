/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/wangyi/GinTemplate/controller/rule"
	"github.com/wangyi/GinTemplate/controller/worker"
	"github.com/wangyi/GinTemplate/dao/mysql"
	"github.com/wangyi/GinTemplate/dao/redis"
	eeor "github.com/wangyi/GinTemplate/error"
	"github.com/wangyi/GinTemplate/model"
	"github.com/wangyi/GinTemplate/tools"
)

func Setup() *gin.Engine {

	r := gin.New()
	r.Use(eeor.ErrHandler())
	r.NoMethod(eeor.HandleNotFound)
	r.NoRoute(eeor.HandleNotFound)
	r.Use(CheckToken())

	r.Static("/static", "./static")
	//用户路由
	user := r.Group("/user")
	{
		//用户登录
		user.POST("/login", worker.Login)
		user.POST("/register", worker.Register)
		user.GET("/generateCaptcha", worker.GenerateCaptcha)
	}

	admin := r.Group("/rule")
	{
		//管理员登录
		admin.POST("/login", rule.Login)
		//获取菜单
		admin.POST("/getMenu", rule.GetMenu)

	}

	r.Run(fmt.Sprintf(":%d", viper.GetInt("app.port")))
	return r
}

//检查检查token
func CheckToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		//检查url 是否在白名单里面
		path := c.Request.URL.Path
		whiteUrl := []string{"/user/register", "/rule/login", "/user/generateCaptcha", "/user/login"}
		if tools.InArray(whiteUrl, path) {
			c.Next()
			return
		}
		//对token 进行检查
		token := c.GetHeader("Token")

		if len(token) == 48 {
			//管理员 token
			//判断token 是否存在
			boolIs, _ := redis.Rdb.HExists("Admin_Token", token).Result()
			if !boolIs {
				tools.JsonWrite(c, -1, []string{}, "非法请求")
				c.Abort()
				return
			}
			_, err := redis.Rdb.Get("Admin_Login_Token_" + token).Result()
			if err != nil {
				tools.JsonWrite(c, -1, []string{}, "登录已经过期")
				c.Abort()
				return
			}
			//设定值  who
			admin := model.AdminModel{}
			err = mysql.DB.Where("token=?", token).First(&admin).Error
			if err != nil {
				tools.JsonWrite(c, -1, []string{}, "系统错误")
				c.Abort()
				return
			}
			c.Set("who", admin)
			c.Next()
			return
		}
		if len(token) == 60 {
			//用户
		}

		tools.JsonWrite(c, -1, []string{}, "非法请求")
		c.Abort()
		return
	}
}
