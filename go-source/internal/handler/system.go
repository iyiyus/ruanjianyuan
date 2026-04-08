package handler

import (
	"go-source/internal/database"
	"go-source/internal/model"
	"go-source/internal/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetUserList 管理员列表（对应前端 /api/user/list）
func GetUserList(c *gin.Context) {
	current, _ := strconv.Atoi(c.DefaultQuery("current", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	var total int64
	database.DB.Model(&model.Admin{}).Count(&total)

	var admins []model.Admin
	database.DB.Offset((current - 1) * size).Limit(size).Find(&admins)

	var records []map[string]interface{}
	for _, a := range admins {
		records = append(records, map[string]interface{}{
			"id":         a.ID,
			"avatar":     a.Avatar,
			"status":     "1",
			"userName":   a.Username,
			"nickName":   a.Nickname,
			"userEmail":  a.Email,
			"userGender": "1",
			"userPhone":  "",
			"userRoles":  []string{"R_SUPER"},
			"createBy":   "system",
			"createTime": time.Unix(a.Createtime, 0).Format("2006-01-02T15:04:05+08:00"),
			"updateBy":   "system",
			"updateTime": time.Unix(a.Updatetime, 0).Format("2006-01-02T15:04:05+08:00"),
		})
	}

	response.Page(c, records, total, current, size)
}

// GetRoleList 角色列表（对应前端 /api/role/list）
func GetRoleList(c *gin.Context) {
	current, _ := strconv.Atoi(c.DefaultQuery("current", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	roles := []map[string]interface{}{
		{
			"roleId":      1,
			"roleName":    "超级管理员",
			"roleCode":    "R_SUPER",
			"description": "拥有所有权限",
			"enabled":     true,
			"createTime":  "2024-01-01T00:00:00+08:00",
		},
	}

	response.Page(c, roles, int64(len(roles)), current, size)
}

// GetAdminLog 管理员日志
func GetAdminLog(c *gin.Context) {
	current, _ := strconv.Atoi(c.DefaultQuery("current", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))

	var total int64
	database.DB.Model(&model.AdminLog{}).Count(&total)

	var list []model.AdminLog
	database.DB.Order("id DESC").Offset((current - 1) * size).Limit(size).Find(&list)

	response.Page(c, list, total, current, size)
}
