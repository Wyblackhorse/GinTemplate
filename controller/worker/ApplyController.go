/**
 * @Author $
 * @Description //TODO $  应用
 * @Date $ $
 * @Param $
 * @return $
 **/
package worker

import (
	"github.com/gin-gonic/gin"
	"github.com/wangyi/GinTemplate/dao/mysql"
	"github.com/wangyi/GinTemplate/model"
)

//获取应用(没有被禁用的)
func GetApply(c *gin.Context) {
	apply := make([]model.Apply, 0)
	err := mysql.DB.Where("status=1").Find(&apply).Error
	if err != nil {
		ReturnErr101(c, "error")
		return
	}
	ReturnSuccessData(c, apply, "success")
	return
}
