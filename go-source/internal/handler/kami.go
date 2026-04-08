package handler

import (
	"crypto/md5"
	"fmt"
	"go-source/internal/database"
	"go-source/internal/model"
	"go-source/internal/response"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func GetKamiList(c *gin.Context) {
	current, _ := strconv.Atoi(c.DefaultQuery("current", "1"))
	size, _ := strconv.Atoi(c.DefaultQuery("size", "20"))
	jh := c.Query("jh")
	kmyp := c.Query("kmyp")

	db := database.DB.Model(&model.Kami{})
	if jh != "" {
		db = db.Where("jh = ?", jh)
	}
	if kmyp != "" {
		db = db.Where("kmyp = ?", kmyp)
	}

	var total int64
	db.Count(&total)

	var list []model.Kami
	db.Order("id DESC").Offset((current - 1) * size).Limit(size).Find(&list)

	response.Page(c, list, total, current, size)
}

type GenKamiReq struct {
	Count      int    `json:"count" binding:"required,min=1,max=1000"`
	Kmyp       int    `json:"kmyp" binding:"required"`
	Prefix     string `json:"prefix"`
	CustomDays int    `json:"customDays"` // kmyp=4时生效
}

func GenerateKami(c *gin.Context) {
	var req GenKamiReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	if req.Kmyp == 4 && req.CustomDays <= 0 {
		response.Fail(c, "请输入自定义天数")
		return
	}

	now := time.Now().Unix()
	gtm := time.Now().Format("20060102150405")
	var codes []string
	var items []model.Kami

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 1; i <= req.Count; i++ {
		code := genKamiCode(req.Prefix, gtm, i, r)
		codes = append(codes, code)
		items = append(items, model.Kami{
			Kami:       code,
			Jh:         0,
			Addtime:    now,
			Usetime:    0,
			Endtime:    0,
			Kmyp:       req.Kmyp,
			CustomDays: req.CustomDays,
		})
	}

	if err := database.DB.Create(&items).Error; err != nil {
		response.Fail(c, "生成失败")
		return
	}

	response.OK(c, gin.H{
		"count": req.Count,
		"codes": strings.Join(codes, "\n"),
	})
}

func DeleteKami(c *gin.Context) {
	id := c.Param("id")
	database.DB.Unscoped().Delete(&model.Kami{}, id)
	response.OKMsg(c, "删除成功", nil)
}

type BatchDeleteKamiReq struct {
	Ids []int `json:"ids" binding:"required,min=1"`
}

func BatchDeleteKami(c *gin.Context) {
	var req BatchDeleteKamiReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	database.DB.Unscoped().Delete(&model.Kami{}, req.Ids)
	response.OKMsg(c, "删除成功", nil)
}

// genKamiCode 与PHP逻辑一致：strtoupper(prefix + substr(md5(gtm+'Km'+i), rand(1,15), 12))
func genKamiCode(prefix, gtm string, i int, r *rand.Rand) string {
	raw := fmt.Sprintf("%sKm%d", gtm, i)
	hash := fmt.Sprintf("%x", md5.Sum([]byte(raw)))
	start := r.Intn(15) + 1 // rand(1,15)
	part := hash[start : start+12]
	return strings.ToUpper(prefix + part)
}
