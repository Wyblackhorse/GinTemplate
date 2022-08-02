package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/wangyi/GinTemplate/dao/mysql"
	"github.com/wangyi/GinTemplate/model"
	"github.com/wangyi/GinTemplate/tools"
	"net/http"
	"strconv"
	"strings"
)

func GetAppUserLoginLog(c *gin.Context) {
	action := c.Query("action")

	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		var total int = 0
		Db := mysql.DB
		fish := make([]model.AppUserLoginLog, 0)
		Db.Table("app_user_login_logs").Count(&total)
		Db = Db.Model(&fish).Offset((page - 1) * limit).Limit(limit)
		if err := Db.Find(&fish).Error; err != nil {
			tools.JsonWrite(c, -101, nil, err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"count":  total,
			"result": fish,
		})
		return
	}

	if action == "Gip" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		var total int = 0
		Db := mysql.DB
		username := c.Query("username")
		//SELECT * FROM `app_user_login_logs` WHERE user_account='clion'  GROUP BY the_login_address
		Db = Db.Where("user_account=?", username)
		fish := make([]model.AppUserLoginLog, 0)
		Db.Table("app_user_login_logs").Count(&total)
		Db = Db.Model(&fish).Offset((page - 1) * limit).Limit(limit).Group("the_login_address")
		if err := Db.Find(&fish).Error; err != nil {
			tools.JsonWrite(c, -101, nil, err.Error())
			return
		}
		type RData struct {
			Ips        string
			Times      string
			Country    string
			IpsNum     int    //同ip数量
			IpsAccount string //同ip的账号数量
		}

		type Country struct {
			Country string
			Nums    int
		}

		var II []Country

		var rr []RData
		for _, v := range fish {
			r := RData{
				Ips:     v.TheLoginAddress,
				Times:   v.AccessTime,
				Country: strings.Split(v.TheLoginSite, "|")[0],
			}
			au := make([]model.AppUserLoginLog, 0)
			mysql.DB.Where("the_login_address=?", v.TheLoginAddress).Group("user_account").Find(&au)

			i := 0
			for _, log := range au {
				if log.UserAccount != username {
					i++
					if r.IpsAccount == "" {
						r.IpsAccount = r.IpsAccount + log.UserAccount
					} else {
						r.IpsAccount = r.IpsAccount + "--" + log.UserAccount
					}
				}
			}
			r.IpsNum = i
			rr = append(rr, r)

			tt := true
			for i2, country := range II {
				if country.Country == r.Country {
					tt = false
					II[i2].Nums = II[i2].Nums + 1
				}
			}

			if tt == true {
				II = append(II, Country{Country: r.Country, Nums: 1})
			}

		}
		type RRT struct {
			A []RData
			B []Country
		}

		var pp RRT
		pp.B=II
		pp.A = rr

		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"count":  len(pp.A),
			"result": pp,
		})
		return

	}

}
