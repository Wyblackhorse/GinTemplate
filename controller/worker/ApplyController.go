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
	"time"
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

	taskDta := make([]model.Task, 0) //hoMap["VipId"]

	if applyIdD, isExist := c.GetQuery("apply_id"); isExist == true {
		//判断应用是否存在
		applyId, _ := strconv.Atoi(applyIdD)
		apply := model.Apply{ID: uint(applyId)}
		if apply.IsExistApply(mysql.DB) == false {
			ReturnErr101(c, "application does not exist")
			return
		}
		mysql.DB.Where("status=? and apply_id=?  AND  task_level <= ?  AND end_time > ?", 1, applyId, whoMap["VipId"],time.Now().Unix()).Find(&taskDta)
	}

	if taskId, isExist := c.GetQuery("task_id"); isExist == true {
		//获取任务
		//fmt.Println(taskId)
		tas := model.Task{}
		err := mysql.DB.Where("id=?", taskId).First(&tas).Error
		if err != nil {
			ReturnErr101(c, "no find  id")
			return
		}
		//ts := model.Task{}
		//mysql.DB.Where("id=?", tas.TaskId).First(&ts)
		//tas.TaskType = ts.TaskType
		//tas.DemandSide = ts.DemandSide
		//tas.TaskUrl = ts.TaskUrl
		ReturnSuccessData(c, tas, "success")
		return
	}

	ReturnData := make([]model.Task, 0)
	for _, i2 := range taskDta {
		err := mysql.DB.Where("task_id=? and worker_id=?  ", i2.ID, whoMap["ID"]).First(&model.TaskOrder{}).Error
		if err != nil {
			ReturnData = append(ReturnData, i2)
		}
	}

	ReturnSuccessData(c, ReturnData, "success")
	return
}

//领取任务
func GetTheApplyTask(c *gin.Context) {
	who, _ := c.Get("who")
	whoMap := who.(map[string]interface{})
	applyId, _ := strconv.Atoi(c.Query("apply_id"))
	taskId, _ := strconv.Atoi(c.Query("task_id"))


	//检查用户的信用分

	//检查任务  应用id 是否存在

	t := model.GetTaskData{WorkerId: int(whoMap["ID"].(uint)), WorkerVipId: whoMap["VipId"].(int), ApplyId: applyId, TaskId: taskId, CreditScore: whoMap["CreditScore"].(int)}
	_, err := t.GetTask(mysql.DB)
	if err != nil {
		ReturnErr101(c, err.Error())
		return
	}
	ReturnSuccess(c, "OK")
	return
}
