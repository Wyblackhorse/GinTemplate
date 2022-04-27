/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package worker

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wangyi/GinTemplate/dao/mysql"
	"github.com/wangyi/GinTemplate/model"
)

//获取  vip  列表
func GetVipList(c *gin.Context) {
	//先获取自己的 vip等级
	who, _ := c.Get("who")

	fmt.Println(who)
	whoMap := who.(map[string]interface{})

	fmt.Println(whoMap)

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


//升级vip
func UpgradeVip(c *gin.Context) {






}
