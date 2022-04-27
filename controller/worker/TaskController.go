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

//获取审核 已完成  已经失败  恶意  已放弃的 任务
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
		//判断任务 id 是否存在
		taskId, _ := strconv.Atoi(c.PostForm("task_id"))
		err2 := mysql.DB.Where("id=?", taskId).First(&model.Task{}).Error
		if err2 != nil {
			ReturnErr101(c, "fail")
			return
		}

		//判断你已经是否已经提交了这个任务
		err2 = mysql.DB.Where("task_id=?", taskId).Where("status=?", 1).Where("worker_id=?", whoMap["ID"]).First(&model.TaskOrder{}).Error
		if err2 == nil {
			ReturnErr101(c, "Don't double submit")
			return
		}

		//判断任务是否超标   获取会员vip等级
		vips := model.Vip{}
		err2 = mysql.DB.Where("vip_ip=?", whoMap["VipId"]).First(&vips).Error
		if err2 != nil {
			ReturnErr101(c, "system is error")
			return
		}
		//获取  今日已经 审核中 和完成的 总数
		var total int
		mysql.DB.Model(&model.TaskOrder{}).Where("status=? or status =?", 1, 3).Where("worker_id=?", whoMap["ID"]).Count(&total)
		if total >= vips.TaskTimes {
			ReturnErr101(c, "The maximum number of tasks exceeded")
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
		nowStr = strconv.Itoa(int(whoMap["ID"].(uint))) + "_" + c.Query("task_id") + nowStr
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
			WorkerId: int(whoMap["ID"].(uint)),
			Created:  time.Now().Unix(),
			ImageUrl: filepath,
			Date:     time.Now().Format("2016_01_02"),
		}
		err2 = mysql.DB.Save(&taskOrder).Error
		if err2 != nil {
			ReturnErr101(c, "err:"+err2.Error())
			return
		}

		ReturnSuccess(c, "success")
		return
	}

}
