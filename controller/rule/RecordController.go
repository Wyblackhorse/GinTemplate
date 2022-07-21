/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $
 **/
package rule

import (
	"encoding/base64"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/wangyi/GinTemplate/dao/mysql"
	"github.com/wangyi/GinTemplate/dao/redis"
	"github.com/wangyi/GinTemplate/model"
	"github.com/wangyi/GinTemplate/tools"
	"net/http"
	"strconv"
	"time"
)

//获取订单
func GetRecords(c *gin.Context) {
	action := c.Query("action")
	//获取订单
	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		kinds := c.Query("kinds")
		role := make([]model.Record, 0)
		Db := mysql.DB
		var total int
		Db = Db.Where("kinds=?", kinds)
		Db.Table("records").Count(&total)
		Db = Db.Model(&model.Record{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
		err := Db.Find(&role).Error
		if err != nil {
			ReturnErr101(c, "ERR:"+err.Error())
			return
		}

		for i, i2 := range role {
			w := model.Worker{}
			err := mysql.DB.Where("id=?", i2.WorkerId).First(&w).Error
			if err == nil {
				role[i].WorkerName = w.Phone
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"count":  total,
			"result": role,
		})
		return
	}
	//提现审核
	if action == "withdrawalAudit" {
		recordId, _ := strconv.Atoi(c.Query("record_id"))
		status, _ := strconv.Atoi(c.Query("status"))
		record := model.Record{ID: uint(recordId)}
		result, _ := record.IsExistRecord(mysql.DB)
		if result == false {
			ReturnErr101(c, "订单不存在")
			return
		}
		if record.Status != 2 || record.Kinds != 2 {
			ReturnErr101(c, "非法修改!")
			return
		}
		//审核 通过
		resultBool, _ := record.WithdrawDeposit(mysql.DB, status)
		if resultBool == false {
			ReturnErr101(c, "审核失败")
			return
		}
		ReturnSuccess(c, "审核成功")
		return

	}

}

func CallBack(c *gin.Context) {

	var aa RecordOrderBack
	err := c.BindJSON(&aa)
	if err != nil {
		ReturnErr101(c, "wrong 1")
		return
	}
	if aa.Code != 200 {
		ReturnErr101(c, "   wrong")
		return
	}

	//base64=> []byte

	decodeString, err1 := base64.StdEncoding.DecodeString(aa.Result.Data)
	if err1 != nil {
		ReturnErr101(c, "sorry  system is  wrong6")
		return
	}

	//解密
	//fmt.Println(aa.Result.Data)
	jsonData, err := tools.RsaDecryptForEveryOne(decodeString)
	if err != nil {
		ReturnErr101(c, "   wrong2"+err.Error())
		return
	}
	var oo RecordOrderBackParameter

	err = json.Unmarshal(jsonData, &oo)
	if err != nil {
		ReturnErr101(c, "   wrong2"+err.Error())
		return
	}
	re := model.Record{}
	PlatformOrderLock, _ := redis.Rdb.SetNX("Record_"+oo.PlatformOrder, time.Now().Unix(), 10*time.Second).Result()
	if PlatformOrderLock == false {
		ReturnErr101(c, "record is doing")
		return
	}
	err = mysql.DB.Where("record_num=?", oo.PlatformOrder).First(&re).Error
	if err != nil {
		ReturnErr101(c, "record is  not exist")
		return
	}
	if re.Status != 2 {
		ReturnErr101(c, "Don't double submit")
		return
	}
	//订单号存在
	wr := model.WorkerBalance{ID: re.WorkerId, AddBalance: oo.AccountPractical, Kinds: 1, Rid: int(re.ID), RecordAccountPractical: oo.AccountPractical}
	_, err = wr.AddBalanceFuc(mysql.DB)
	if err != nil {
		ReturnErr101(c, err.Error())
		return
	}
	ReturnSuccess(c, "OK")
	return
}
