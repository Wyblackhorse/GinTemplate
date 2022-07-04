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
	"github.com/wangyi/GinTemplate/dao/mysql"
	"github.com/wangyi/GinTemplate/model"
	"strconv"
)

//获取  vip  列表
func GetVipList(c *gin.Context) {
	//先获取自己的 vip等级
	who, _ := c.Get("who")
	whoMap := who.(map[string]interface{})
	vips := make([]model.Vip, 0)
	err := mysql.DB.Find(&vips).Error
	if err != nil {
		ReturnErr101(c, "fail")
		return
	}
	//
	dataOne := VipsReturnData{
		MeVip: whoMap["VipId"],
		Its:   vips,
	}
	//
	//fmt.Println(dataOne)
	ReturnSuccessData(c, dataOne, "success")
	return
}

//升级vip  购买vip
func UpgradeVip(c *gin.Context) {
	who, _ := c.Get("who")
	whoMap := who.(map[string]interface{})
	vipId, _ := strconv.Atoi(c.Query("vip_id"))
	v := model.Vip{ID: uint(vipId)}
	//判断vip 是否存在
	returnA, price := v.ReturnVipPrice(mysql.DB)
	if returnA == false {
		ReturnErr101(c, "Sorry, purchase failed")
		return
	}
	//获取用户的余额
	w := model.WorkerBalance{ID: int(whoMap["ID"].(uint)), AddBalance: price, Kinds: 4, VipID: vipId}
	_, err := w.AddBalanceFuc(mysql.DB)
	if err != nil {
		ReturnErr101(c, err.Error())
		return
	}

	d := model.DailyStatistics{TodayAddVipNums: 1}
	d.SetEverydayData(mysql.DB)
	ReturnSuccess(c, "success")
	return
}
