package controller

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/wangyi/GinTemplate/dao/mysql"
	"github.com/wangyi/GinTemplate/model"
	"github.com/wangyi/GinTemplate/tools"
	"github.com/xuri/excelize/v2"
	"strconv"
	"strings"
)

func UploadFiles(c *gin.Context) {
	action := c.Query("action")
	//传的 recharge 表
	file, err := c.FormFile("file")
	if err != nil {
		tools.JsonWrite(c, -101, nil, err.Error())
		return
	}

	fileType := strings.Split(file.Filename, ".")
	if fileType[1] != "xlsx" {
		tools.JsonWrite(c, -101, nil, "上传格式不对")
		return
	}

	err = c.SaveUploadedFile(file, file.Filename)
	if err != nil {
		tools.JsonWrite(c, -101, nil, "上传错误:"+
			err.Error())
		return
	}

	f, err2 := excelize.OpenFile(file.Filename)
	if err2 != nil {
		tools.JsonWrite(c, -101, nil, "上传错误:"+
			err2.Error())
		return
	}

	defer func() {
		if err := f.Close(); err != nil {
			fmt.Println(err)
		}
	}()

	if action == "recharge" {
		rows, err1 := f.GetRows("recharge")
		if err1 != nil {
			tools.JsonWrite(c, -101, nil, "上传错误:"+err1.Error())
			return
		}
		go func() {
			for k, row := range rows {
				//fmt.Println(value)
				if k == 0 {
					continue //break
				}
				add := model.Recharge{}
				add.Created = row[0]
				add.Updated = row[1]
				add.Remark = row[2]
				add.UserId, _ = strconv.Atoi(row[4])
				add.Username = row[5]
				add.Recharge, _ = strconv.ParseFloat(row[6], 64)
				add.CloseAccount, _ = strconv.ParseFloat(row[7], 64)
				add.TopUpAward, _ = strconv.ParseFloat(row[8], 64)
				add.Status = row[9]
				add.Kinds, _ = strconv.Atoi(row[10])
				add.Serial = row[11]
				add.ThreeOrders = row[12]
				add.TopUpChannel, _ = strconv.Atoi(row[13])
				add.Classify = row[14]
				err1 := mysql.DB.Where("serial=?", add.Serial).First(&model.Recharge{}).Error
				if err1 != nil {
					mysql.DB.Save(&add)
				}
			}
		}()
	}







	tools.JsonWrite(c, 200, nil, "OK")

}
