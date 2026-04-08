package handler

import (
	"go-source/internal/database"
	"go-source/internal/model"
	"go-source/internal/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetCategoryList(c *gin.Context) {
	current, _ := strconv.Atoi(c.DefaultQuery("current", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	name := c.Query("name")
	status := c.Query("status")

	db := database.DB.Model(&model.Category{})
	if name != "" {
		db = db.Where("name LIKE ?", "%"+name+"%")
	}
	if status != "" {
		db = db.Where("status = ?", status)
	}

	var total int64
	db.Count(&total)

	var list []model.Category
	db.Order("weigh DESC").Offset((current - 1) * size).Limit(size).Find(&list)

	response.Page(c, list, total, current, size)
}

func GetCategory(c *gin.Context) {
	id := c.Param("id")
	var item model.Category
	if err := database.DB.First(&item, id).Error; err != nil {
		response.Fail(c, "记录不存在")
		return
	}
	response.OK(c, item)
}

func CreateCategory(c *gin.Context) {
	var item model.Category
	if err := c.ShouldBindJSON(&item); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	now := time.Now().Unix()
	item.Createtime = now
	item.Updatetime = now
	if item.Status == "" {
		item.Status = "normal"
	}
	if err := database.DB.Create(&item).Error; err != nil {
		response.Fail(c, "创建失败")
		return
	}
	response.OK(c, item)
}

func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var item model.Category
	if err := database.DB.First(&item, id).Error; err != nil {
		response.Fail(c, "记录不存在")
		return
	}
	if err := c.ShouldBindJSON(&item); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	item.Updatetime = time.Now().Unix()
	database.DB.Save(&item)
	response.OK(c, item)
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	if err := database.DB.Unscoped().Delete(&model.Category{}, id).Error; err != nil {
		response.Fail(c, "删除失败")
		return
	}
	response.OKMsg(c, "删除成功", nil)
}
