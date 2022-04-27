/**
 * @Author $
 * @Description //TODO $
 * @Date $ $
 * @Param $
 * @return $  权限管理
 **/
package rule

import (
	"github.com/gin-gonic/gin"
	"github.com/wangyi/GinTemplate/dao/mysql"
	"github.com/wangyi/GinTemplate/model"
	"github.com/wangyi/GinTemplate/tools"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//获取角色列表
func GetRole(c *gin.Context) {
	action := c.Query("action")
	//获取基本数据
	if action == "GET" {
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		role := make([]model.Role, 0)
		Db := mysql.DB
		var total int
		Db = Db.Model(&model.Role{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
		Db.Table("roles").Count(&total)
		err := Db.Find(&role).Error
		if err != nil {
			ReturnErr101(c, "ERR:"+err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"count":  total,
			"result": role,
		})
		return
	}
	//新增角色
	if action == "ADD" {
		name := c.Query("Name")
		err := mysql.DB.Where("name=?", name).First(&model.Role{}).Error
		if err == nil {
			ReturnErr101(c, "不可以重复添加")
			return
		}
		mysql.DB.Save(&model.Role{Name: name, Status: 1, Created: time.Now().Unix()})
		ReturnSuccess(c, "OK")
		return
	}
	//更新数据  不包含权限
	if action == "UPDATE" {
		id := c.Query("ID")
		update := model.Role{}
		//修改名字
		if name, isExist := c.GetQuery("Name"); isExist == true {
			update.Name = name
		}
		//修改状态    // 1 正常 2封禁
		if a, isExits := c.GetQuery("Status"); isExits == true {
			if id == "1" {
				ReturnErr101(c, "超级管理员不可以修改")
				return
			}
			status, _ := strconv.Atoi(a)
			update.Status = status
		}
		err := mysql.DB.Model(&model.Role{}).Where("id=?", id).Update(&update).Error
		if err != nil {
			ReturnErr101(c, "修改失败:"+err.Error())
			return
		}
		ReturnSuccess(c, "success")
		return
	}
	//获取权限
	if action == "AuthorityList" {
		//获取权限表
		//获取管理员id
		id := c.Query("ID")

		//判断 角色是否存在
		err2 := mysql.DB.Where("id=?", id).First(&model.Role{}).Error
		if err2 != nil {
			ReturnErr101(c, "fail")
			return
		}
		menus := make([]model.Menu, 0)
		err := mysql.DB.Raw("SELECT *  FROM menus ").Scan(&menus).Error
		if err != nil {
			ReturnErr101(c, "err:"+err.Error())
			return
		}
		var m MenusJsonData
		for _, v := range menus {
			err := mysql.DB.Where("role_id=? and menu_id=?", id, v.ID).First(&model.RoleAndMenu{}).Error
			v.IfChoose = 2
			if err == nil {
				//找到了
				v.IfChoose = 1
			}
			if v.Level == 0 {
				var a MenusJson
				a.Top = v
				m.Data = append(m.Data, a)
			} else {
				m.Data[v.Level-1].Sec = append(m.Data[v.Level-1].Sec, v)
			}
		}
		ReturnSuccessData(c, m, "OK")
		return
	}
	//修改权限
	if action == "UpAuthorityList" {
		//管理员id
		RoleIdOne := c.Query("ID")
		//权限id  用@分割
		if delIdString, isExist := c.GetQuery("delIds"); isExist == true {
			idsArray := strings.Split(delIdString, "@")
			for _, v := range idsArray {
				mysql.DB.Model(&model.RoleAndMenu{}).Delete("menu_id=? and role_id", v, RoleIdOne)
			}
		}
		//新增权限
		RoleIdOneId, _ := strconv.Atoi(RoleIdOne)
		if delIdString, isExist := c.GetQuery("addIds"); isExist == true {
			idsArray := strings.Split(delIdString, "@")
			for _, v := range idsArray {
				id, _ := strconv.Atoi(v)
				mysql.DB.Model(&model.RoleAndMenu{}).Save(&model.RoleAndMenu{RoleId: RoleIdOneId, MenuId: id, Created: time.Now().Unix()})
			}
		}
		ReturnSuccess(c, "OK")
		return
	}
	ReturnSuccess(c, "no action")
	return
}

//获取权限列表
func GetJurisdiction(c *gin.Context) {

	action := c.Query("action")
	//获取权限列表
	if action == "GET" {

		admin := make([]model.AdminModel, 0)
		page, _ := strconv.Atoi(c.Query("page"))
		limit, _ := strconv.Atoi(c.Query("limit"))
		Db := mysql.DB
		var total int
		Db = Db.Model(&model.AdminModel{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
		Db.Table("admin_models").Count(&total)
		err := Db.Find(&admin).Error
		if err != nil {
			ReturnErr101(c, "ERR:"+err.Error())
			return
		}
		for k, v := range admin {
			roles := model.Role{}
			err := mysql.DB.Where("id=?", v.RoleId).First(&roles).Error
			if err == nil {
				admin[k].RoleName = roles.Name
			}
		}
		c.JSON(http.StatusOK, gin.H{
			"code":   1,
			"count":  total,
			"result": admin,
		})
		return
	}
	//添加数据
	if action == "ADD" {
		username := c.Query("username")
		token := tools.RandStringRunes(48)
		roleId := c.Query("role_id")
		roleIdI, _ := strconv.Atoi(roleId)
		password := c.Query("password")
		add := model.AdminModel{RoleId: roleIdI, Created: time.Now().Unix(), Username: username, Token: token, Password: password}
		//判断是是否重复添加
		err := mysql.DB.Where("username=?", username).First(&model.AdminModel{}).Error
		if err == nil {
			ReturnErr101(c, "不要重复添加")
			return
		}
		err = mysql.DB.Save(&add).Error
		if err != nil {
			ReturnErr101(c, "error:"+err.Error())
			return
		}
		ReturnSuccess(c, "OK")
		return

	}
	//更新数据
	if action == "UPDATE" {
		id := c.Query("id")
		//判断数据是否存在
		err := mysql.DB.Where("id=?", id).First(&model.AdminModel{}).Error
		if err != nil {
			ReturnErr101(c, "fail")
			return
		}
		//更新状态 单独的按钮
		if status, isExist := c.GetQuery("status"); isExist == true {
			statusInt, _ := strconv.Atoi(status)
			mysql.DB.Model(&model.AdminModel{}).Where("id=?", id).Update(&model.AdminModel{Status: statusInt})
			ReturnSuccess(c, "OK")
			return
		}

		//更新其他
		username := c.Query("username")
		roleId := c.Query("role_id")
		roleIdI, _ := strconv.Atoi(roleId)
		password := c.Query("password")
		mysql.DB.Model(&model.AdminModel{}).Where("id=?", id).Update(&model.AdminModel{Username: username, RoleId: roleIdI, Password: password})
		ReturnSuccess(c, "OK")
		return
	}

	ReturnSuccess(c, "no action")
	return

}



