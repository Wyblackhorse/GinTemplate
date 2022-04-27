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
	"time"
)

//设置任务
func SetTasks(c *gin.Context) {
	action := c.Query("action")
	//获取基本数据
	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		role := make([]model.Task, 0)
		Db := mysql.DB
		var total int
		Db = Db.Model(&model.Role{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
		Db.Table("tasks").Count(&total)
		err := Db.Find(&role).Error
		if err != nil {
			ReturnErr101(c, "ERR:"+err.Error())
			return
		}

		for k, v := range role {
			apply := model.Apply{}
			err := mysql.DB.Where("id=?", v.ApplyId).First(&apply).Error
			if err == nil {
				role[k].ApplyName = apply.Name
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"count":  total,
			"result": role,
		})
		return
	}
	//添加操作
	if action == "ADD" {
		applyID, _ := strconv.Atoi(c.Query("apply_id"))
		remark := c.Query("remark")
		applyType, _ := strconv.Atoi(c.Query("task_type"))
		endTime, _ := strconv.ParseInt(c.Query("end_time"), 10, 64)
		price, _ := strconv.ParseFloat(c.Query("price"), 64)
		taskNum, _ := strconv.Atoi(c.Query("task_num"))
		TaskLevel, _ := strconv.Atoi(c.Query("task_level"))
		taskUrl := c.Query("task_url")
		add := model.Task{ApplyId: applyID, Remark: remark, TaskType: applyType, EndTime: endTime, Price: price, TaskNum: taskNum, TaskUrl: taskUrl, Created: time.Now().Unix(), Status: 1, TaskLevel: TaskLevel}
		err := mysql.DB.Save(&add).Error
		if err != nil {
			ReturnErr101(c, "添加失败")
			return
		}
		ReturnSuccess(c, "添加成功")
		return
	}
	//更新操作
	if action == "UPDATE" {

	}

}

//对任务进行审核  SetTasksOrder
func SetTasksOrder(c *gin.Context) {

	//who, _ := c.Get("who")
	//whoMap := who.(map[string]interface{})

	action := c.Query("action")
	//获取任务   状态    2审核中 3已完成 4以失败 5恶意 6已放弃
	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		role := make([]model.TaskOrder, 0)
		Db := mysql.DB
		var total int
		Db = Db.Model(&model.TaskOrder{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
		Db.Table("task_orders").Count(&total)
		err := Db.Find(&role).Error
		if err != nil {
			ReturnErr101(c, "ERR:"+err.Error())
			return
		}
		for k, v := range role {
			apply := model.Task{}
			err := mysql.DB.Where("id=?", v.TaskId).First(&apply).Error
			if err == nil {
				role[k].TaskName = apply.Remark
			}
			worker := model.Worker{}
			err = mysql.DB.Where("id=?", v.WorkerId).First(&worker).Error
			if err == nil {
				role[k].WorkerName = worker.Phone
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"count":  total,
			"result": role,
		})
		return
	}
	//更新审核
	if action == "UPDATE" {
		//获取订单id
		orderId := c.Query("order_id")
		workerId, _ := strconv.Atoi(c.Query("worker_id"))
		//判断订单id 是否存在
		order := model.TaskOrder{}
		err := mysql.DB.Where("id=?", orderId).Where("worker_id=?", workerId).First(&order).Error
		if err != nil {
			ReturnErr101(c, "订单不存在")
			return
		}
		status, _ := strconv.Atoi(c.Query("status"))
		if status == order.Status {
			//状态相等
			ReturnErr101(c, "不要重复提交")
			return
		}

		if status == 3 {
			//已完成   加钱  获取这单任务的 单价
			task := model.Task{}
			err := mysql.DB.Where("id=?", order.TaskId).First(&task).Error
			if err != nil {
				ReturnErr101(c, "任务id非法")
				return
			}
			WorkerBalance := model.WorkerBalance{ID: workerId, Kinds: 3, AddBalance: task.Price, OrderId: int(order.ID)}
			WorkerBalanceBool, err := WorkerBalance.AddBalanceFuc(mysql.DB)
			if WorkerBalanceBool == false {
				ReturnErr101(c, "审核失败")
				return
			}
			ReturnSuccess(c, "审核成功")
			return
		} else if status == 4 {
			//失败
			err := mysql.DB.Model(&model.TaskOrder{}).Where("id=?", orderId).Update(&model.TaskOrder{Status: 5, Updated: time.Now().Unix()}).Error
			if err != nil {
				ReturnErr101(c, "审核失败")
				return
			}
			ReturnSuccess(c, "审核成功")
			return
		} else if status == 5 {
			//恶意的  扣除信用分
			err := mysql.DB.Model(&model.TaskOrder{}).Where("id=?", orderId).Update(&model.TaskOrder{Status: 5, Updated: time.Now().Unix()}).Error
			if err != nil {
				ReturnErr101(c, "审核失败")
				return
			}
			worker := model.Worker{ID: uint(workerId)}
			worker.SubtractCreditScore(mysql.DB)
			ReturnSuccess(c, "审核成功")
			return

		}

	}

	ReturnErr101(c, "错误的action")
	return

}
