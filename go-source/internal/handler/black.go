package handler

import (
	"go-source/internal/database"
	"go-source/internal/model"
	"go-source/internal/response"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func GetBlackList(c *gin.Context) {
	current, _ := strconv.Atoi(c.DefaultQuery("current", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	udid := c.Query("udid")

	db := database.DB.Model(&model.Black{})
	if udid != "" {
		db = db.Where("udid LIKE ?", "%"+udid+"%")
	}

	var total int64
	db.Count(&total)

	var list []model.Black
	db.Order("id DESC").Offset((current - 1) * size).Limit(size).Find(&list)

	response.Page(c, list, total, current, size)
}

type CreateBlackReq struct {
	Udid   string `json:"udid" binding:"required"`
	Reason string `json:"reason"`
}

func CreateBlack(c *gin.Context) {
	var req CreateBlackReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	var exist model.Black
	if database.DB.Where("udid = ?", req.Udid).First(&exist).Error == nil {
		response.Fail(c, "该UDID已在黑名单中")
		return
	}
	item := model.Black{
		Udid:    req.Udid,
		Reason:  req.Reason,
		Addtime: time.Now().Unix(),
	}
	database.DB.Create(&item)
	response.OK(c, item)
}

func DeleteBlack(c *gin.Context) {
	id := c.Param("id")
	database.DB.Unscoped().Delete(&model.Black{}, id)
	response.OKMsg(c, "删除成功", nil)
}

// GetBlacklistPublic 公开黑名单接口，无需认证，返回所有黑名单UDID及原因
func GetBlacklistPublic(c *gin.Context) {
	var list []model.Black
	database.DB.Order("id DESC").Find(&list)

	type Item struct {
		Udid    string `json:"udid"`
		Reason  string `json:"reason"`
		Addtime int64  `json:"addtime"`
	}
	items := make([]Item, 0, len(list))
	for _, v := range list {
		if v.Udid != "" {
			items = append(items, Item{
				Udid:    v.Udid,
				Reason:  v.Reason,
				Addtime: v.Addtime,
			})
		}
	}
	response.OK(c, gin.H{
		"total": len(items),
		"list":  items,
	})
}

// CheckBlacklist 查询单个UDID是否在黑名单，无需认证
func CheckBlacklist(c *gin.Context) {
	udid := c.Query("udid")
	if udid == "" {
		response.Fail(c, "udid不能为空")
		return
	}
	var black model.Black
	err := database.DB.Where("udid = ?", udid).First(&black).Error
	if err != nil {
		response.OK(c, gin.H{"blocked": false, "reason": ""})
		return
	}
	response.OK(c, gin.H{"blocked": true, "reason": black.Reason})
}
