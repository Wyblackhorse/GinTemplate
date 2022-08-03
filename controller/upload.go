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
	"time"
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
				add.Date = strings.Split(add.Updated, " ")[0]
				loc, _ := time.LoadLocation("Asia/Shanghai")                        //设置时区
				tt, _ := time.ParseInLocation("2006-01-02", add.Date, loc) //2006-01-02 15:04:05是转换的格式如php的"Y-m-d H:i:s"
				add.DateTimestamp = tt.Unix()

				err1 := mysql.DB.Where("serial=?", add.Serial).First(&model.Recharge{}).Error
				if err1 != nil {
					mysql.DB.Save(&add)
				}
			}
		}()
	}

	if action == "walletRecord" {
		rows, err1 := f.GetRows("walletrecord")
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
				add := model.WalletRecord{}
				add.Created = row[0]
				add.Updated = row[1]
				add.Note = row[2]
				add.SerialNumber, _ = strconv.Atoi(row[4])
				add.UserName = row[5]
				add.Amount, _ = strconv.ParseFloat(row[6], 64)
				add.Kinds = row[7]
				add.BeforeTheValue, _ = strconv.Atoi(row[8])
				add.AfterTheValue, _ = strconv.Atoi(row[9])
				add.Serial = row[10]
				err1 := mysql.DB.Where("serial=?", add.Serial).First(&model.WalletRecord{}).Error
				if err1 != nil {
					mysql.DB.Save(&add)
				}
			}
		}()
	}

	if action == "withdraw" {
		rows, err1 := f.GetRows("withdraw")
		if err1 != nil {
			tools.JsonWrite(c, -101, nil, "上传错误:"+err1.Error())
			return
		}
		go func() {
			for k, row := range rows {
				if k == 0 {
					continue
				}

				add := model.Withdraw{}
				add.RecordId, _ = strconv.Atoi(row[0])
				add.TheUserId, _ = strconv.Atoi(row[1])
				add.TheUserName = row[2]
				add.WithdrawalAmount, _ = strconv.ParseFloat(row[3], 64)
				add.Status = row[4]
				add.Kinds, _ = strconv.Atoi(row[5])
				add.TheActualAmount, _ = strconv.ParseFloat(row[6], 64)
				add.TheSettlementAmount, _ = strconv.ParseFloat(row[7], 64)
				add.Poundage, _ = strconv.ParseFloat(row[8], 64)
				add.Rate, _ = strconv.ParseFloat(row[9], 64)
				add.OrderNo = row[10]
				add.Classification = row[11]
				add.ChannelID, _ = strconv.Atoi(row[12])

				//add.ThirdPartyTrackingNumber = row[13]
				err1 := mysql.DB.Where("record_id=?", add.RecordId).First(&model.Withdraw{}).Error
				if err1 != nil {
					mysql.DB.Save(&add)
				}
			}
		}()
	}

	if action == "appUser" {
		rows, err1 := f.GetRows("appuser")
		if err1 != nil {
			tools.JsonWrite(c, -101, nil, "上传错误:"+err1.Error())
			return
		}
		go func() {
			for k, row := range rows {
				if k == 0 {
					continue
				}
				add := model.AppUser{}
				add.UserNumber, _ = strconv.Atoi(row[0])
				add.TheHigherTheID, _ = strconv.Atoi(row[1])
				add.UpperLayerUserName = row[2]
				add.GeneralAgentID, _ = strconv.Atoi(row[3])
				add.TheGeneralAgentOf = row[4]
				add.Username = row[5]
				add.MobilePhoneNo = row[6]
				add.UserMailbox = row[7]
				add.State = row[8]
				add.RegistrationTime = row[9]
				add.RegisteredIP = row[10]
				add.Updated = row[11]
				add.InviteCode = row[12]
				add.LastLoginTime = row[13]
				add.RealName = row[14]
				add.TestNo = row[15]
				add.Grouping = row[16]

				err1 := mysql.DB.Where("username=?", add.Username).First(&model.AppUser{}).Error
				if err1 != nil {
					mysql.DB.Save(&add)
				}
			}
		}()
	}

	if action == "appUserLoginLog" {
		rows, err1 := f.GetRows("appuser_login_log")
		if err1 != nil {
			tools.JsonWrite(c, -101, nil, "上传错误:"+err1.Error())
			return
		}

		go func() {
			for k, row := range rows {
				if k == 0 {
					continue
				}
				add := model.AppUserLoginLog{}
				add.Create = row[0]
				add.Update = row[1]
				add.Remark = row[2]
				add.SerialNumber, _ = strconv.Atoi(row[3])
				add.UserID, _ = strconv.Atoi(row[4])
				add.UserAccount = row[5]
				add.LoginID = row[6]
				add.TheLoginAddress = row[7]
				add.TheLoginSite = row[8]
				add.Browser = row[9]
				add.OperatingSystem = row[10]
				add.OperatingInformation = row[11]
				add.AccessTime = row[12]
				mysql.DB.Save(&add)
			}
		}()

	}

	if action == "bettingRecord" {
		rows, err1 := f.GetRows("bettingRecord")
		if err1 != nil {
			tools.JsonWrite(c, -101, nil, "上传错误:"+err1.Error())
			return
		}
		go func() {
			for k, row := range rows {
				if k == 0 {
					continue
				}
				add := model.BettingRecord{}
				add.UserID, _ = strconv.Atoi(row[0])
				add.UserName = row[1]
				add.EventID, _ = strconv.Atoi(row[2])
				add.State = row[3]
				add.BreakEven = row[4]
				add.BetAmount, _ = strconv.ParseFloat(row[5], 64)
				add.Payout, _ = strconv.ParseFloat(row[6], 64)
				add.ProfitAndLoss = row[7]
				add.OrderID, _ = strconv.Atoi(row[8])
				add.FullHalf = row[9]
				add.PositiveWaveCounter = row[10]
				add.Score = row[11]
				add.Alliance = row[12]
				add.TopVSBottom = row[13]
				add.Yield, _ = strconv.ParseFloat(row[14], 64)
				add.StringOfNo = row[15]
				add.TheFinalResult = row[16]
				add.TheStartTime = row[17]
				add.BetOnTime = row[18]
				add.Poundage, _ = strconv.ParseFloat(row[19], 64)
				add.PoundagePer, _ = strconv.ParseFloat(row[20], 64)
				add.BetDate = strings.Split(add.BetOnTime, " ")[0]
				mysql.DB.Save(&add)

			}
		}()

	}

	tools.JsonWrite(c, 200, nil, "OK")

}
