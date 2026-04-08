package handler

import (
	"go-source/internal/database"
	"go-source/internal/model"
	"go-source/internal/response"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetMonitorList(c *gin.Context) {
	current, _ := strconv.Atoi(c.DefaultQuery("current", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	udid := c.Query("udid")
	identity := c.Query("identity")

	db := database.DB.Model(&model.Monitor{})
	if udid != "" {
		db = db.Where("udid LIKE ?", "%"+udid+"%")
	}
	if identity != "" {
		db = db.Where("identity = ?", identity)
	}

	var total int64
	db.Count(&total)

	var list []model.Monitor
	db.Order("id DESC").Offset((current - 1) * size).Limit(size).Find(&list)

	response.Page(c, list, total, current, size)
}

func DeleteMonitor(c *gin.Context) {
	id := c.Param("id")
	database.DB.Unscoped().Delete(&model.Monitor{}, id)
	response.OKMsg(c, "删除成功", nil)
}

func ClearMonitor(c *gin.Context) {
	database.DB.Where("1 = 1").Delete(&model.Monitor{})
	response.OKMsg(c, "清空成功", nil)
}
