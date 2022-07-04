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

//获取会员等级  修改 删除
func GetVipLevel(c *gin.Context) {
	action := c.Query("action")
	//获取vip
	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		role := make([]model.Vip, 0)
		Db := mysql.DB
		var total int
		Db.Table("vips").Count(&total)
		Db = Db.Model(&model.Vip{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
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
	//修改
	if action == "UPDATE" {
		id := c.Query("ID")
		err := mysql.DB.Where("id=?", id).First(&model.Vip{}).Error
		if err != nil {
			ReturnErr101(c, "id不存在")
			return
		}

		ups := make(map[string]interface{})

		//名字
		if name, _isExist := c.GetQuery("Name"); _isExist == true {
			ups["Name"] = name
		}

		//等级
		if name, _isExist := c.GetQuery("Level"); _isExist == true {
			ups["Level"], _ = strconv.Atoi(name)
		}

		//价格
		if name, _isExist := c.GetQuery("Price"); _isExist == true {
			ups["Price"], _ = strconv.ParseFloat(name, 64)
		}

		//每单收入
		if name, _isExist := c.GetQuery("Account"); _isExist == true {
			ups["Account"], _ = strconv.ParseFloat(name, 64)
		}
		//次数
		if name, _isExist := c.GetQuery("TaskTimes"); _isExist == true {
			ups["TaskTimes"], _ = strconv.Atoi(name)
		}

		err = mysql.DB.Model(model.Vip{}).Where("id=?", id).Update(ups).Error
		if err != nil {
			ReturnErr101(c, "更新失败")
			return
		}
		ReturnSuccess(c, "更新成功")
		return
	}
	//添加
	if action == "ADD" {

	}
	//删除
	if action == "DEL" {

	}
	ReturnErr101(c, "no action")
	return

}

//获取会员成员
func GetVipWorkers(c *gin.Context) {

	action := c.Query("action")
	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		role := make([]model.Worker, 0)
		Db := mysql.DB
		var total int

		//vip等级
		if vipLevel, isExits := c.GetQuery("vip_id"); isExits == true {
			vipId, _ := strconv.Atoi(vipLevel)
			Db = Db.Where("vip_id=?", vipId)
		}
		//账号模糊搜索
		if phone, isExist := c.GetQuery("phone"); isExist == true {
			Db = Db.Where("phone LIKE ?", "%"+phone+"%")
		}

		//注册时间的范围
		if start, isExist := c.GetQuery("start"); isExist == true {
			if end, isExist := c.GetQuery("end"); isExist == true {
				Db = Db.Where("created BETWEEN ? AND ?", start, end)
			}
		}

		//余额范围
		if start, isExist := c.GetQuery("minBalance"); isExist == true {
			if end, isExist := c.GetQuery("maxBalance"); isExist == true {
				Db = Db.Where("balance BETWEEN ? AND ?", start, end)
			}
		}
		Db.Table("workers").Count(&total)
		Db = Db.Model(&model.Worker{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
		err := Db.Find(&role).Error
		if err != nil {
			ReturnErr101(c, "ERR:"+err.Error())
			return
		}
		for k, v := range role {
			//判断会员等级
			vip := model.Vip{ID: uint(v.VipId)}
			role[k].VipName = vip.GetLevelName(mysql.DB)
		}
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"count":  total,
			"result": role,
		})
		return
	}

	//更新用户
	if action == "UPDATE" {
		workerId, _ := strconv.Atoi(c.Query("worker_id"))
		worker := model.Worker{ID: uint(workerId)}
		if !worker.IsExist(mysql.DB) {
			ReturnErr101(c, "这个用户不存在")
			return
		}
		update := map[string]interface{}{}
		//需改状态    状态 1限制  2良好 3优秀  4封号
		if updateData, isExist := c.GetQuery("status"); isExist == true {
			update["Status"], _ = strconv.Atoi(updateData)
		}

		//修改信用分
		if updateData, isExist := c.GetQuery("CreditScore"); isExist == true {
			update["CreditScore"], _ = strconv.Atoi(updateData)
		}
		//修改会员等级
		if updateData, isExist := c.GetQuery("VipId"); isExist == true {
			//判断这个vip 等级是否存在
			VipId, _ := strconv.Atoi(updateData)
			vip := model.Vip{ID: uint(VipId)}
			b, _ := vip.ReturnVipPrice(mysql.DB)
			if b == false {
				ReturnErr101(c, "vipId 不存在")
				return
			}
			update["VipId"] = VipId
		}

		//真实姓名
		if updateData, isExist := c.GetQuery("Username"); isExist == true {
			update["Username"] = updateData
		}
		//登录密码
		if updateData, isExist := c.GetQuery("Password"); isExist == true {
			update["Password"] = updateData
		}
		//取款密码
		if updateData, isExist := c.GetQuery("PayPassword"); isExist == true {
			update["PayPassword"] = updateData
		}
		//电子邮箱
		if updateData, isExist := c.GetQuery("EMail"); isExist == true {
			update["EMail"] = updateData
		}
		//会员到期时间   (0 无限期)
		if updateData, isExist := c.GetQuery("VipExpire"); isExist == true {
			update["VipExpire"], _ = strconv.ParseInt(updateData, 10, 0)
		}

		err := mysql.DB.Model(model.Worker{}).Where("id=?", workerId).Update(update).Error
		if err != nil {
			ReturnErr101(c, "修改失败")
			return
		}
		ReturnSuccess(c, "修改成功")
		return

	}

	ReturnErr101(c, "no action")
	return
}

//会员银行

//管理操作资金  (管理员操作用户的余额)
func ChangeMoneyForAdmin(c *gin.Context) {
	workerId, _ := strconv.Atoi(c.Query("worker_id"))
	w := model.Worker{ID: uint(workerId)}
	if w.IsExist(mysql.DB) == false {
		ReturnErr101(c, "用户不存在")
		return
	}

	AddBalance, _ := strconv.ParseFloat(c.Query("money"), 64)
	kinds, _ := strconv.Atoi(c.Query("kinds"))

	balance := model.WorkerBalance{ID: workerId, AddBalance: AddBalance, Kinds: kinds, IfAdmin: true, Remark: c.Query("remark")}
	_, err := balance.AddBalanceFuc(mysql.DB)
	if err != nil {
		ReturnErr101(c, err.Error())
		return
	}
	ReturnSuccess(c, "OK")
	return

}
