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
	"log"
	"net/http"
	"strings"
)

func Setup() *gin.Engine {

	r := gin.New()
	r.Use(Cors())
	r.Use(eeor.ErrHandler())
	r.NoMethod(eeor.HandleNotFound)
	r.NoRoute(eeor.HandleNotFound)
	r.Static("/static", "./static")
	r.Use(CheckToken())
	//用户路由
	user := r.Group("/user")
	{
		//用户登录 注册 获取验证嘛
		user.POST("/login", worker.Login)
		user.POST("/register", worker.Register)
		user.GET("/generateCaptcha", worker.GenerateCaptcha)

		//获取应用列表
		user.GET("/apply/getApply", worker.GetApply)
		//GetApplyTask 获取应用的任务
		user.GET("/apply/getApplyTask", worker.GetApplyTask)
		//GetTheApplyTask
		user.GET("/apply/getTheApplyTask", worker.GetTheApplyTask)

		//获取 vips GetVipList
		user.GET("/vip/getVipList", worker.GetVipList)
		//购买vip UpgradeVip
		user.GET("/vip/upgradeVip", worker.UpgradeVip)

		//GetTask
		user.GET("/task/getTask", worker.GetTask)
		//GetTaskOrder   //获取订单任务     提交审核
		user.GET("/task/getTaskOrder", worker.GetTaskOrder)
		user.POST("/task/getTaskOrder", worker.GetTaskOrder)
		//提现 RecordController
		user.GET("/record/getRecord", worker.GetRecord)

		//获取余额宝
		user.GET("/moneyManagement/shoppingMoneyManagement", worker.ShoppingMoneyManagement)
		//获取个人信息 Information
		user.GET("/getInformation", worker.Information)
		//团队报表
		user.GET("team/teamTasks", worker.TeamTasks)
		//DailyReport
		user.GET("DailyReport", worker.DailyReport)
		//购买云管家 GetCloudHousekeeper
		user.GET("GetCloudHousekeeper", worker.GetCloudHousekeeper)



	}

	admin := r.Group("/rule")
	{
		//首页
		admin.GET("/homePage", rule.HomePage)
		//管理员登录
		admin.POST("/login", rule.Login)

		//用于支付成功时候的回调
		admin.POST("/callBack", rule.CallBack)

		//获取权限列表 		//获取菜单
		admin.GET("/jurisdictionManagement/roleManagement/getRole", rule.GetRole)
		admin.GET("/jurisdictionManagement/roleManagement/getJurisdiction", rule.GetJurisdiction)
		//获取应用列表  GetApplyList
		admin.GET("/applyManagement/applyList/getApplyList", rule.GetApplyList)
		//任务管理 taskManagement SetTasks taskList     审核任务 reviewTheTask   SetCollection
		admin.GET("/taskManagement/taskList/setTasks", rule.SetTasks)
		admin.GET("/taskManagement/reviewTheTask/setTasksOrder", rule.SetTasksOrder)
		admin.GET("/taskManagement/collection/setCollection", rule.SetCollection)
		//系统设置  SetConfig
		admin.GET("/settingManagement/basicSetting/setConfig", rule.SetConfig)
		//LanternSlide  幻灯片设置
		admin.GET("/settingManagement/lanternSlide/lanternSlide", rule.LanternSlide)
		admin.POST("/settingManagement/lanternSlide/lanternSlide", rule.LanternSlide)
		//会员管理   普通会员  会员等级 		//GetBill(余额变动)  ChangeMoneyForAdmin(修改用户余额)  获取银行
		admin.GET("/memberManagement/gradeOfMembership/getVipLevel", rule.GetVipLevel)
		admin.GET("/memberManagement/regularMembers/getVipWorkers", rule.GetVipWorkers)
		admin.GET("/memberManagement/gradeOfMembership/getBill", rule.GetBill)
		admin.GET("/memberManagement/gradeOfMembership/changeMoneyForAdmin", rule.ChangeMoneyForAdmin)
		admin.GET("/memberManagement/memberBank/getBank", rule.GetBank)
		//账单管理  充值账单  提现账单  佣金账单 推广奖励
		admin.GET("/billManagement/withdrawalBill/getRecord", rule.GetRecords)
		//日志管理  GetLoggerList adminLog(管理操作日志)
		admin.GET("/logManagement/adminLog/getAdminLog", rule.GetAdminLog)
		//获取管理者(日志管理)
		admin.GET("/logManagement/adminLog/getAdmin", rule.GetAdmin)
		//获取系统日志  GetSystemLog
		admin.GET("/logManagement/systemLog/getSystemLog", rule.GetSystemLog)
		//报表 数据  GetStatementEveryday(每日报表)
		admin.GET("/statementManagement/statementEveryday/getStatementEveryday", rule.GetStatementEveryday)
		//团队报表  GetTeamStatistics
		admin.GET("/statementManagement/statementTeam/getTeamStatistics", rule.GetTeamStatistics)
		//全局统计   GetGlobalStatistics
		admin.GET("/statementManagement/globalStatistics/getGlobalStatistics", rule.GetGlobalStatistics)
		//余额宝管理  balanceManagement
		admin.GET("/balanceManagement/productList/getYuEBaoList", rule.GetYuEBaoList)
		// 获取余额购买记录   GetYuEBaoPurchaseHistory
		admin.GET("/balanceManagement/purchaseHistory/getYuEBaoPurchaseHistory", rule.GetYuEBaoPurchaseHistory)
		//CloudHousekeeperDoTask 计划任务
		admin.GET("CloudHousekeeperDoTask", rule.CloudHousekeeperDoTask)
	}

	r.Run(fmt.Sprintf(":%d", viper.GetInt("app.port")))
	return r
}

//检查检查token
func CheckToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		//检查url 是否在白名单里面
		path := c.Request.URL.Path
		whiteUrl := []string{"/user/register", "/rule/login", "/user/generateCaptcha", "/user/login", "/rule/callBack"}
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

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin") //请求头部
		if origin != "" {
			//接收客户端发送的origin （重要！）
			c.Writer.Header().Set("Access-Control-Allow-Origin", origin)
			//服务器支持的所有跨域请求的方法
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE,UPDATE")
			//允许跨域设置可以返回其他子段，可以自定义字段
			c.Header("Access-Control-Allow-Headers", "Authorization, Content-Length, X-CSRF-Token, Token,session")
			// 允许浏览器（客户端）可以解析的头部 （重要）
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers")
			//设置缓存时间
			c.Header("Access-Control-Max-Age", "172800")
			//允许客户端传递校验信息比如 cookie (重要)
			c.Header("Access-Control-Allow-Credentials", "true")
		}

		//允许类型校验
		if method == "OPTIONS" {
			c.JSON(http.StatusOK, "ok!")
		}

		defer func() {
			if err := recover(); err != nil {
				log.Printf("Panic info is: %v", err)
			}
		}()

		c.Next()
	}
}
