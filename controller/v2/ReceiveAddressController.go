package V2

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/wangyi/GinTemplate/dao/mysql"
	"github.com/wangyi/GinTemplate/model"
	"github.com/wangyi/GinTemplate/tools"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// GetReceiveAddress 获取地址管理
func GetReceiveAddress(c *gin.Context) {
	action := c.Query("action")
	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		role := make([]model.ReceiveAddress, 0)
		Db := mysql.DB
		var total int
		Db = Db.Model(&model.ReceiveAddress{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
		Db.Table("receive_addresses").Count(&total)
		err := Db.Find(&role).Error
		if err != nil {
			tools.ReturnError101(c, "ERR:"+err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":  0,
			"count": total,
			"data":  role,
		})
		return
	}

	if action == "ADD" {
		result := model.CreateNewReceiveAddress(mysql.DB, viper.GetString("eth.ThreeUrl"))
		if !result {
			tools.ReturnError101(c, "添加失败")
			return
		}
		tools.ReturnError200(c, "添加成功")
		return
	}

}

// Collection 资金归集
func Collection(c *gin.Context) {
	req := make(map[string]interface{})
	req["gas"] = c.Query("gas")
	req["min"] = c.Query("min")
	if req["gas"] == "" || req["min"] == "" {
		tools.ReturnError101(c, "非法参数")
		return
	}

	if addr, isExits := c.GetQuery("addr"); isExits == true {
		if addr != "" {
			addArray := strings.Split(addr, "@")
			req["addrs"] = addArray

		}

	}

	req["ts"] = time.Now().UnixMilli()
	_, err := tools.HttpRequest(viper.GetString("eth.ThreeUrl")+"/collect", req, viper.GetString("eth.ApiKey"))
	if err != nil {
		tools.ReturnError101(c, "归集失败")
		return
	}
	tools.ReturnError200(c, "归集成功")
	return
}

// GetAllMoney 获取总余额
func GetAllMoney(c *gin.Context) {
	rec := make([]model.ReceiveAddress, 0)
	err := mysql.DB.Find(&rec).Error
	if err != nil {
		tools.ReturnError101(c, "获取失败")
		return
	}
	var sumMoney float64
	for _, v := range rec {
		sumMoney = sumMoney + v.Money
	}
	tools.ReturnError200Data(c, sumMoney, "获取成功")
	return
}
