package handler

import (
	"go-source/internal/database"
	"go-source/internal/model"
	"go-source/internal/response"

	"github.com/gin-gonic/gin"
)

func GetConfigList(c *gin.Context) {
	var list []model.Config
	database.DB.Order("id ASC").Find(&list)
	response.OK(c, list)
}

func GetConfigMap(c *gin.Context) {
	var list []model.Config
	database.DB.Find(&list)
	m := make(map[string]string)
	for _, v := range list {
		m[v.Name] = v.Value
	}
	response.OK(c, m)
}

type UpdateConfigReq struct {
	Configs map[string]string `json:"configs" binding:"required"`
}

func UpdateConfig(c *gin.Context) {
	var req UpdateConfigReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	for name, value := range req.Configs {
		database.DB.Model(&model.Config{}).Where("name = ?", name).Update("value", value)
	}
	response.OKMsg(c, "保存成功", nil)
}
