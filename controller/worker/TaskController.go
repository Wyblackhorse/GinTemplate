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
	"strings"
	"time"
)

//获取任务
func GetTask(c *gin.Context) {
	who, _ := c.Get("who")
	whoMap := who.(map[string]interface{})

	tasks := make([]model.Task, 0)
	//只显示 状态没有 为 1 的 订单
	err := mysql.DB.Where("status=?", 1).Find(&tasks).Error
	if err != nil {
		ReturnErr101(c, "fail")
		return
	}

	for k, v := range tasks {
		apply := model.Apply{}
		err = mysql.DB.Where("id=?", v.ApplyId).First(&apply).Error
		if err == nil {
			tasks[k].ApplyName = apply.Name
		}
		//查询这个玩家师父提交了任务   status=1 正常运行的
		err := mysql.DB.Where("task_id=?", v.ID).Where("worker_id=?", whoMap["ID"]).Where("status=?", 1).First(&model.TaskOrder{}).Error
		if err != nil {
			tasks[k].WorkerStatus = 1
		} else {
			tasks[k].WorkerStatus = 2
		}

	}
	ReturnSuccessData(c, tasks, "success")
	return

}

//获取审核 已完成  已经失败  恶意  已放弃的 任务  进行中
func GetTaskOrder(c *gin.Context) {
	action := c.Query("action")
	who, _ := c.Get("who")
	whoMap := who.(map[string]interface{})
	//获取数据 
	if action == "GET" {
		status, _ := strconv.Atoi(c.Query("status"))
		taskOrder := make([]model.TaskOrder, 0)
		err := mysql.DB.Where("worker_id=?", whoMap["ID"]).Where("status=?", status).Find(&taskOrder).Error
		if err != nil {
			ReturnErr101(c, "fail")
			return
		}
		for k, v := range taskOrder {
			task := model.Task{}
			err := mysql.DB.Where("id=?", v.TaskId).First(&task).Error
			if err == nil {
				apply := model.Apply{}
				err = mysql.DB.Where("id=?", task.ApplyId).First(&apply).Error
				if err == nil {
					taskOrder[k].TaskName = apply.Name
				}
			}
		}
		ReturnSuccessData(c, taskOrder, "success")
		return
	}
	//提交审核
	if action == "Submit" {
		//判断任务 id 是否存在  并且还要判断 这个任务是否还有效
		taskId, _ := strconv.Atoi(c.PostForm("task_order_id"))
		err2 := mysql.DB.Where("id=?", taskId).Where("status=?", 1).First(&model.TaskOrder{}).Error
		if err2 != nil {
			ReturnErr101(c, "fail")
			return
		}
		//获取 图片
		file, err := c.FormFile("file")
		if err != nil {
			ReturnErr101(c, "fail")
			return
		}
		if file.Size > 67444 {
			ReturnErr101(c, "Picture is too big")
			return
		}
		//判断是否是图片
		nameArray := strings.Split(file.Filename, ".")
		if len(nameArray) != 2 {
			ReturnErr101(c, "fail")
			return
		}
		if nameArray[1] != "png" {
			ReturnErr101(c, "It must be a PNG image")
			return
		}
		nowStr := time.Now().Format("20060102150405")
		nowStr = strconv.Itoa(int(whoMap["ID"].(uint))) + "_" + c.Query("task_order_id") + nowStr
		filepath := "./static/upload/" + nowStr + ".png"
		err = c.SaveUploadedFile(file, filepath)
		if err != nil {
			ReturnErr101(c, "err:"+err.Error())
			return
		}
		//上传成功 生成订单
		taskOrder := model.TaskOrder{
			Status:   2, //审核中 给管理员审核
			TaskId:   taskId,
			ImageUrl: filepath,
			Updated:  time.Now().Unix(),
		}

		err2 = mysql.DB.Model(&model.TaskOrder{}).Where("id=?", taskId).Update(&taskOrder).Error
		if err2 != nil {
			ReturnErr101(c, "err:"+err2.Error())
			return
		}
		d := model.DailyStatistics{TodayAddVipNums: 1}
		d.SetEverydayData(mysql.DB)
		ReturnSuccess(c, "success")
		return
	}
}