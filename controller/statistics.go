package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wangyi/GinTemplate/dao/mysql"
	"github.com/wangyi/GinTemplate/model"
	"github.com/wangyi/GinTemplate/tools"
	"net/http"
	"strconv"
	"time"
)

// SetStatistics 生成数据
func SetStatistics(c *gin.Context) {
	action := c.Query("action")

	//统计数据的初始化 生成
	if action == "INIT" {
		//获取充值人数
		st := make([]model.Recharge, 0)
		//按照日期分组(目的:取出充值的所有日期)
		mysql.DB.Group("date").Find(&st)

		fmt.Println("初始化开始.---")

		for _, recharge := range st {
			add := model.Statistics{}
			//获取充值的日期   SELECT   *   FROM  recharges   WHERE date='2022-06-14'    GROUP BY user_id
			//充值人数
			mysql.DB.Model(&model.Recharge{}).Where("date=? and status=?", recharge.Date, "成功").Group("user_id").Count(&add.TodayWithdrawNums)
			st1 := make([]model.Recharge, 0)
			mysql.DB.Where("date=? and status=?", recharge.Date, "成功").Find(&st1)
			var peopleNum []int
			for _, i2 := range st1 {
				if tools.InArray(peopleNum, i2.UserId) == false {
					peopleNum = append(peopleNum, recharge.UserId)
					//今日首冲人数
					loc, _ := time.LoadLocation("Asia/Shanghai")                    //设置时区
					tt, _ := time.ParseInLocation("2006-01-02", recharge.Date, loc) //2006-01-02 15:04:05是转换的格式如php的"Y-m-d H:i:s"
					err := mysql.DB.Where("date_timestamp < ? and status=? and user_id=?", tt.Unix(), "成功", i2.UserId).First(&model.Recharge{}).Error
					if err != nil {
						add.TodayWithdrawFirst = add.TodayWithdrawFirst + 1
					}
					//今日充值没有下注的人数
					err = mysql.DB.Where("bet_date =?  and user_id=?", recharge.Date, i2.UserId).First(&model.BettingRecord{}).Error
					if err == nil { //找到了下注
						add.TodayWithdrawBetNums = add.TodayWithdrawBetNums + 1 //今日充值下注活跃人数+1
					}
					//今日存款提现数据
					//mysql.DB.Model(&model.Withdraw{}).Where("status=?   and  date =? ","通过", recharge.Date).
					mysql.DB.Raw("select sum(withdrawal_amount) as  today_withdraw_deposit from  withdraws where date  =? and status=?", recharge.Date, "通过").Scan(&add)

				}
			}
			//今日私自下注人数
			mysql.DB.Model(&model.BettingRecord{}).Where("bet_date=? and break_even=? and state=?", recharge.Date, "正常", "已结算").Group("user_id").Count(&add.TodayPrivatelyBetNums)

			//今日总的下注人数
			mysql.DB.Model(&model.BettingRecord{}).Where("bet_date=?  and state=?", recharge.Date, "已结算").Group("user_id").Count(&add.TodayAllBetNums)
			add.TodayWithdrawNoBetNums = add.TodayWithdrawNums - add.TodayWithdrawBetNums
			//判断统计数据是否已经存在这个日期的数据了
			stt := model.Statistics{}
			err := mysql.DB.Where("date=?", recharge.Date).First(&stt).Error
			if err != nil {
				//没有数据插入
				add.Updated = time.Now().Unix()
				add.Created = time.Now().Unix()
				add.Date = recharge.Date
				mysql.DB.Save(&add)

			} else {
				//更新数
				add.Updated = time.Now().Unix()
				mysql.DB.Model(&model.Statistics{}).Where("id=?", stt.ID).Update(&add)
			}

			//fmt.Println(add)

		}
		fmt.Println("初始化结束.---")

		tools.JsonWrite(c, 200, nil, "OK")
		return
	}
	//获取数据
	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))

		var total int = 0
		Db := mysql.DB
		fish := make([]model.Statistics, 0)
		Db.Table("statistics").Count(&total)
		Db = Db.Model(&fish).Offset((page - 1) * limit).Limit(limit).Order("date desc")

		if err := Db.Find(&fish).Error; err != nil {
			tools.JsonWrite(c, -101, nil, err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"count":  total,
			"result": fish,
		})

		return
	}
	//首页
	if action == "HOMEPAGE" {
		st := model.Statistics{}
		mysql.DB.Where("date=?", time.Now().Format("2006-01-02")).First(&st)
		mysql.DB.Raw("SELECT sum(today_withdraw_no_bet_nums) as in_all_withdraw_no_bet_nums ,sum(today_withdraw_bet_nums) as in_all_withdraw_no_bet_nums , sum(today_withdraw_bet_nums) as in_all_withdraw_bet_nums    FROM statistics ").Scan(&st)
		mysql.DB.Raw("SELECT sum(today_withdraw_nums) as in_all_withdraw_nums ,sum(today_withdraw_first) as in_all_withdraw_first , sum(today_privately_bet_nums) as in_all_privately_bet_nums    FROM statistics ").Scan(&st)
		mysql.DB.Raw("SELECT sum(today_withdraw_deposit) as in_all_withdraw_deposit ,sum(today_commission) as in_all_commission   FROM statistics ").Scan(&st)
		tools.JsonWrite(c, 200, st, "OK")
		return
	}

	//获取单独日期的详情
	if action == "DETAIL" {
		date := c.Query("date")
		var peopleNum []int
		//获取首冲会员的名单
		//获取充值的日期   SELECT   *   FROM  recharges   WHERE date='2022-06-14'    GROUP BY user_id
		//充值人数
		st1 := make([]model.Recharge, 0)
		mysql.DB.Where("date=? and status=?", date, "成功").Find(&st1)
		returnData := model.Detail{}
		for _, i2 := range st1 {
			if tools.InArray(peopleNum, i2.UserId) == false {
				APP := model.AppUser{}
				APPEAR := mysql.DB.Where("user_number=?", i2.UserId).First(&APP).Error
				peopleNum = append(peopleNum, i2.UserId)
				//今日首冲人数
				loc, _ := time.LoadLocation("Asia/Shanghai")           //设置时区
				tt, _ := time.ParseInLocation("2006-01-02", date, loc) //2006-01-02 15:04:05是转换的格式如php的"Y-m-d H:i:s"
				err := mysql.DB.Where("date_timestamp < ? and status=? and user_id=?", tt.Unix(), "成功", i2.UserId).First(&model.Recharge{}).Error
				if err != nil {
					//判定为今日首次充值
					if APPEAR == nil {
						returnData.TodayWithdrawFirstDetail = append(returnData.TodayWithdrawFirstDetail, model.User{Username: APP.Username, TheGeneralAgentOf: APP.TheGeneralAgentOf})
					}
				}

				//今日充值没有下注的人数
				err = mysql.DB.Where("bet_date =?  and user_id=?", date, i2.UserId).First(&model.BettingRecord{}).Error
				if err != nil { //没有找到这个投注记录
					if APPEAR == nil {
						fmt.Println("----")
						returnData.TodayWithdrawNoBetNumsDetail = append(returnData.TodayWithdrawNoBetNumsDetail, model.User{Username: APP.Username, TheGeneralAgentOf: APP.TheGeneralAgentOf})
					}
				} else {
					//今日下注人员明细
					if APPEAR == nil {

						returnData.TodayWithdrawBetNumsDetail = append(returnData.TodayWithdrawBetNumsDetail, model.User{Username: APP.Username, TheGeneralAgentOf: APP.TheGeneralAgentOf})
					}

				}
				//今日私自下注人数
				err = mysql.DB.Where("bet_date=? and break_even!=?", date, "正常").First(&model.BettingRecord{}).Error
				if err == nil {
					if APPEAR == nil {
						returnData.TodayPrivatelyBetNumsDetail = append(returnData.TodayPrivatelyBetNumsDetail, model.User{Username: APP.Username, TheGeneralAgentOf: APP.TheGeneralAgentOf})
					}
				}

			}
		}
		tools.JsonWrite(c, 200, returnData, "OK")
		return
	}

}
