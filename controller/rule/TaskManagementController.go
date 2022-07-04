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
	//获取基本数据 (任务列表)
	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		role := make([]model.Task, 0)
		Db := mysql.DB

		if applyID, isExist := c.GetQuery("apply_id"); isExist == true {
			Db = Db.Where("apply_id=?", applyID)

		}
		if applyID, isExist := c.GetQuery("remark"); isExist == true {
			Db = Db.Where("remark=?", applyID)

		}
		if applyID, isExist := c.GetQuery("task_type"); isExist == true {
			Db = Db.Where("task_type=?", applyID)

		}
		if start, isExist := c.GetQuery("start"); isExist == true {
			if applyID, isExist := c.GetQuery("end"); isExist == true {
				Db = Db.Where("created  between  ? and ?", start, applyID)
			}
		}
		if applyID, isExist := c.GetQuery("status"); isExist == true {
			Db = Db.Where("status=?", applyID)
		}
		if applyID, isExist := c.GetQuery("id"); isExist == true {
			Db = Db.Where("id=?", applyID)
		}
		if applyID, isExist := c.GetQuery("task_url"); isExist == true {
			Db = Db.Where("task_url=?", applyID)
		}

		var total int
		Db.Table("tasks").Count(&total)
		Db = Db.Model(&model.Role{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
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
		demandSide := c.Query("demand_side")
		taskUrl := c.Query("task_url")
		if taskUrl == "" {
			ReturnErr101(c, "taskUrl is  not  ''")
			return
		}

		//不能添加同一条链接
		err2 := mysql.DB.Where("task_url=?", taskUrl).First(&model.Task{}).Error
		if err2 == nil {
			ReturnErr101(c, "不要重复添加")
			return
		}

		add := model.Task{ApplyId: applyID, Remark: remark, TaskType: applyType, EndTime: endTime, Price: price, TaskNum: taskNum, TaskUrl: taskUrl, Created: time.Now().Unix(), Status: 1, TaskLevel: TaskLevel, AlreadyGet: taskNum, DemandSide: demandSide}
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
		id := c.Query("id")
		err := mysql.DB.Where("id=?", id).First(&model.Task{}).Error
		if err != nil {
			ReturnErr101(c, "Illegal to submit")
			return
		}
		ups := make(map[string]interface{})
		if applyID, isExist := c.GetQuery("apply_id"); isExist == true {
			ups["ApplyID"], _ = strconv.Atoi(applyID)

		}
		if applyID, isExist := c.GetQuery("remark"); isExist == true {
			ups["Remark"] = applyID

		}
		if applyID, isExist := c.GetQuery("task_type"); isExist == true {
			ups["TaskType"], _ = strconv.Atoi(applyID)
		}
		if applyID, isExist := c.GetQuery("end_time"); isExist == true {
			ups["EndTime"], _ = strconv.ParseInt(applyID, 10, 64)
		}
		if applyID, isExist := c.GetQuery("price"); isExist == true {
			ups["Price"], _ = strconv.ParseFloat(applyID, 64)

		}
		if applyID, isExist := c.GetQuery("task_num"); isExist == true {
			ups["TaskNum"], _ = strconv.Atoi(applyID)

		}
		if applyID, isExist := c.GetQuery("task_level"); isExist == true {
			ups["TaskLevel"], _ = strconv.Atoi(applyID)

		}
		if applyID, isExist := c.GetQuery("demand_side"); isExist == true {
			ups["DemandSide"] = applyID

		}
		if applyID, isExist := c.GetQuery("task_url"); isExist == true {
			ups["TaskUrl"] = applyID

		}
		err = mysql.DB.Model(&model.Task{}).Where("id=?", id).Update(ups).Error
		if err != nil {
			ReturnErr101(c, "更新失败")
			return
		}
		ReturnSuccess(c, "更新成功")
		return
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
		status, _ := strconv.Atoi(c.Query("status"))
		role := make([]model.TaskOrder, 0)
		Db := mysql.DB.Where("status=?", status)
		var total int
		Db.Table("task_orders").Count(&total)
		Db = Db.Model(&model.TaskOrder{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
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
				role[k].TaskType = apply.TaskType
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
			err := mysql.DB.Model(&model.TaskOrder{}).Where("id=?", orderId).Update(&model.TaskOrder{Status: 4, Updated: time.Now().Unix()}).Error
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

		} else if status == 6 {
			//撤销

			err := mysql.DB.Model(&model.TaskOrder{}).Where("id=?", orderId).Update(&model.TaskOrder{Status: 6, Updated: time.Now().Unix()}).Error
			if err != nil {
				ReturnErr101(c, "撤销失败")
				return
			}
			data := model.GetTaskData{Status: 6, TaskId: order.TaskId}
			_, _ = data.GetTask(mysql.DB)
			ReturnSuccess(c, "撤销成功")
		}
		return
	}
	ReturnErr101(c, "错误的action")
	return
}

//添加采集任务
func SetCollection(c *gin.Context) {
	action := c.Query("action")
	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		role := make([]model.Collection, 0)
		Db := mysql.DB
		var total int
		Db.Table("collections").Count(&total)

		Db = Db.Model(&model.Collection{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
		Db.Find(&role)

		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"count":  total,
			"result": role,
		})
		return
	}

	if action == "ADD" {
		TaskType, _ := strconv.Atoi(c.Query("task_type"))
		Expiry, _ := strconv.Atoi(c.Query("expiry"))
		TaskNum, _ := strconv.Atoi(c.Query("task_num"))
		TaskLevel, _ := strconv.Atoi(c.Query("task_level"))
		kinds, _ := strconv.Atoi(c.Query("kinds"))
		Price, _ := strconv.ParseFloat(c.Query("price"), 64)
		cc := model.Collection{
			TaskUrl:    c.Query("url"),
			TaskType:   TaskType,
			Expiry:     Expiry,
			TaskNum:    TaskNum,
			Price:      Price,
			DemandSide: c.Query("demand_side"),
			TaskLevel:  TaskLevel,
			Remark:     c.Query("remark"),
			Kinds:      kinds,
		}
		_, err := cc.Add(mysql.DB)
		if err != nil {
			ReturnErr101(c, err.Error())
			return
		}
		ReturnSuccess(c, "添加成功")
		return
	}

}
