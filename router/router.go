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
	"github.com/fatih/structs"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/wangyi/GinTemplate/controller/rule"
	"github.com/wangyi/GinTemplate/controller/worker"
	"github.com/wangyi/GinTemplate/dao/mysql"
	"github.com/wangyi/GinTemplate/dao/redis"
	eeor "github.com/wangyi/GinTemplate/error"
	"github.com/wangyi/GinTemplate/model"
	"github.com/wangyi/GinTemplate/tools"
	"strings"
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
		//用户登录 注册 获取验证嘛
		user.POST("/login", worker.Login)
		user.POST("/register", worker.Register)
		user.GET("/generateCaptcha", worker.GenerateCaptcha)

		//获取应用列表
		user.GET("/apply/getApply", worker.GetApply)
		//获取 vips GetVipList
		user.GET("/vip/getVipList", worker.GetVipList)
		//GetTask
		user.GET("/task/getTask", worker.GetTask)
		//GetTaskOrder   //获取订单任务     提交审核
		user.GET("/task/getTaskOrder", worker.GetTaskOrder)
		user.POST("/task/getTaskOrder", worker.GetTaskOrder)
		//提现 RecordController
		user.GET("/record/getRecord", worker.GetRecord)



	}

	admin := r.Group("/rule")
	{
		//管理员登录
		admin.POST("/login", rule.Login)
		//获取权限列表 		//获取菜单
		admin.GET("/jurisdictionManagement/roleManagement/getRole", rule.GetRole)
		admin.GET("/jurisdictionManagement/roleManagement/getJurisdiction", rule.GetJurisdiction)
		//获取应用列表  GetApplyList
		admin.GET("/applyManagement/applyList/getApplyList", rule.GetApplyList)
		//任务管理 taskManagement SetTasks taskList     审核任务 reviewTheTask
		admin.GET("/taskManagement/taskList/setTasks", rule.SetTasks)
		admin.GET("/taskManagement/reviewTheTask/setTasksOrder", rule.SetTasksOrder)
		//系统设置  SetConfig
		admin.GET("/settingManagement/basicSetting/setConfig", rule.SetConfig)
		//会员管理   普通会员  会员等级
		admin.GET("/memberManagement/gradeOfMembership/getVipLevel", rule.GetVipLevel)
		admin.GET("/memberManagement/regularMembers/getVipWorkers", rule.GetVipWorkers)
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

			m3 := structs.Map(&admin)
			c.Set("who", m3)
			//c.Set("who", admin)
			c.Next()
			return
		}
		if len(token) == 60 {
			//用户

			//校验路径权限
			path := c.Request.URL.Path
			pathArray := strings.Split(path, "/")
			if pathArray[1] != "user" {
				tools.JsonWrite(c, -1, []string{}, "非法访问")
				c.Abort()
				return
			}

			//查看token是否非法
			boolIs, _ := redis.Rdb.HExists("Worker_Token", token).Result()
			if !boolIs {
				tools.JsonWrite(c, -1, []string{}, "非法请求")
				c.Abort()
				return
			}
			_, err := redis.Rdb.Get("Worker_Login_Token_" + token).Result()
			if err != nil {
				tools.JsonWrite(c, -1, []string{}, "登录已经过期")
				c.Abort()
				return
			}

			//设定值  who
			admin := model.Worker{}
			err = mysql.DB.Where("token=?", token).First(&admin).Error
			if err != nil {
				tools.JsonWrite(c, -1, []string{}, "系统错误")
				c.Abort()
				return
			}
			m3 := structs.Map(&admin)
			c.Set("who", m3)
			c.Next()

			return

		}

		tools.JsonWrite(c, -1, []string{}, "非法请求")
		c.Abort()
		return
	}
}
