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
	"github.com/wangyi/GinTemplate/model"
	"net/http"
	"strconv"
)

//获取余额包 活动列表
func GetYuEBaoList(c *gin.Context) {
	action := c.Query("action")
	//获取活动数据
	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		role := make([]model.YuEBao, 0)
		Db := mysql.DB
		var total int
		Db.Table("yu_e_baos").Count(&total)
		Db = Db.Model(&model.YuEBao{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
		err := Db.Find(&role).Error
		if err != nil {
			ReturnErr101(c, "ERR:"+err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"count":  total,
			"result": role,
		})
		return
	}
	//添加数据
	if action == "ADD" {
		name := c.Query("name")
		minMoney, _ := strconv.ParseFloat(c.Query("min_money"), 64)
		interestRate, _ := strconv.ParseFloat(c.Query("interest_rate"), 64)
		Days, _ := strconv.Atoi(c.Query("days"))
		bao := model.YuEBao{Name: name, MinMoney: minMoney, InterestRate: interestRate, Days: Days}
		result, _ := bao.AddYuEBao(mysql.DB)
		if result == false {
			ReturnErr101(c, "添加失败")
			return
		}
		ReturnSuccess(c, "添加成功")
		return
	}

	if action == "UPDATE" {
		id := c.Query("id")
		//判断是否存在
		err := mysql.DB.Where("id=?", id).First(&model.YuEBao{}).Error
		if err != nil {
			ReturnErr101(c, "非法请求")
			return
		}

		ups := make(map[string]interface{}, 0)
		//名字
		if name, isE := c.GetQuery("name"); isE == true {
			ups["Name"] = name
		}

		//最小的购买金额
		if name, isE := c.GetQuery("min_money"); isE == true {
			minMoney, _ := strconv.ParseFloat(name, 64)
			ups["MinMoney"] = minMoney
		}

		//利率
		if name, isE := c.GetQuery("interest_rate"); isE == true {
			minMoney, _ := strconv.ParseFloat(name, 64)
			ups["InterestRate"] = minMoney
		}

		//时间
		if name, isE := c.GetQuery("days"); isE == true {
			Days, _ := strconv.Atoi(name)
			ups["InterestRate"] = Days
		}

		//状态
		if name, isE := c.GetQuery("status"); isE == true {
			Days, _ := strconv.Atoi(name)
			ups["Status"] = Days
		}

		mysql.DB.Model(&model.YuEBao{}).Where("id=?", id).Update(ups)

		ReturnSuccess(c, "修改成功")
		return
	}

	ReturnErr101(c, "no action")
	return
}

//获取余额宝的购买记录
func GetYuEBaoPurchaseHistory(c *gin.Context) {
	action := c.Query("action")
	//获取余额宝购买记录
	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		role := make([]model.MoneyManagement, 0)
		Db := mysql.DB
		var total int
		Db.Table("money_managements").Count(&total)
		Db = Db.Model(&model.MoneyManagement{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
		err := Db.Find(&role).Error
		if err != nil {
			ReturnErr101(c, "ERR:"+err.Error())
			return
		}

		for i, i2 := range role {
			//获取用户名名字
			w := model.Worker{}
			err := mysql.DB.Where("id=?", i2.WorkerId).First(&w).Error
			if err == nil {
				role[i].WorkerName = w.Phone
			}
			//获取产品名字
			p := model.YuEBao{}
			err = mysql.DB.Where("id=?", i2.YuEBaoId).First(&p).Error
			if err == nil {
				role[i].YuEBaoName = p.Name
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"count":  total,
			"result": role,
		})
		return
	}

}
