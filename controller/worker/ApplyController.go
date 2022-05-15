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
	"strconv"
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

//获取没人领取的任务
func GetApplyTask(c *gin.Context) {
	who, _ := c.Get("who")
	whoMap := who.(map[string]interface{})
	applyId, _ := strconv.Atoi(c.Query("apply_id"))
	apply := model.Apply{ID: uint(applyId)}
	//判断应用是否存在
	if apply.IsExistApply(mysql.DB) == false {
		ReturnErr101(c, "application does not exist")
		return
	}
	//获取任务
	taskDta := make([]model.Task, 0) //hoMap["VipId"]
	err := mysql.DB.Where("status=? and apply_id=?  AND  task_level <= ?", 1, applyId, whoMap["VipId"]).Find(&taskDta).Error
	if err != nil {
		ReturnSuccessData(c, taskDta, "success")
		return
	}
	ReturnSuccessData(c, taskDta, "success")
	return
}

//领取任务
func GetTheApplyTask(c *gin.Context) {
	who, _ := c.Get("who")
	whoMap := who.(map[string]interface{})
	applyId, _ := strconv.Atoi(c.Query("apply_id"))
	taskId, _ := strconv.Atoi(c.Query("task_id"))
	t := model.GetTaskData{WorkerId: int(whoMap["ID"].(uint)), WorkerVipId: whoMap["VipId"].(int), ApplyId: applyId, TaskId: taskId}
	_, err := t.GetTask(mysql.DB)
	if err != nil {
		ReturnErr101(c, err.Error())
		return
	}
	ReturnSuccess(c, "OK")
	return
}
