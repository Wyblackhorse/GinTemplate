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
)

//团队报表
func TeamTasks(c *gin.Context) {
	action := c.Query("action")
	who, _ := c.Get("who")

	whoMap := who.(map[string]interface{})

	//团队报表
	if action == "task" {
		worker := model.Worker{ID: whoMap["ID"].(uint)}
		data := worker.GetTeamStatement(mysql.DB)
		ReturnSuccessData(c, data, "success")
		return
	}
	//我的团队
	if action == "myTeam" {
		worker := model.Worker{ID: 17}
		data := worker.GetMyTeamInformation(mysql.DB)
		ReturnSuccessData(c, data, "success")
		return
	}

}
