package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-source/internal/database"
	"go-source/internal/model"
	"go-source/internal/response"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type FetchSourceReq struct {
	URL  string `json:"url" binding:"required"`
	UDID string `json:"udid"`
}

// FetchSource 请求远程软件源数据
func FetchSource(c *gin.Context) {
	var req FetchSourceReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	url := req.URL
	if req.UDID != "" {
		if strings.Contains(url, "?") {
			url += "&udid=" + req.UDID
		} else {
			url += "?udid=" + req.UDID
		}
	}

	client := &http.Client{Timeout: 30 * time.Second}
	httpReq, _ := http.NewRequest("GET", url, nil)
	httpReq.Header.Set("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 14_7_1 like Mac OS X) AppleWebKit/605.1.15")
	httpReq.Header.Set("Accept", "application/json, text/plain, */*")

	resp, err := client.Do(httpReq)
	if err != nil {
		response.Fail(c, "请求失败: "+err.Error())
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)

	var data map[string]interface{}
	if err := json.Unmarshal(body, &data); err != nil {
		response.Fail(c, "无效的软件源数据格式")
		return
	}

	// 处理加密格式 {"appstore": "base64..."}
	if b64, ok := data["appstore"].(string); ok {
		decoded, err := base64.StdEncoding.DecodeString(b64)
		if err != nil {
			response.Fail(c, "解密失败")
			return
		}
		var inner map[string]interface{}
		if err := json.Unmarshal(decoded, &inner); err != nil {
			response.Fail(c, "解析加密数据失败")
			return
		}
		data = inner
	}

	if _, ok := data["apps"]; !ok {
		response.Fail(c, "软件源数据中没有apps字段")
		return
	}

	response.OK(c, data)
}

type ImportAppsReq struct {
	Apps []map[string]interface{} `json:"apps" binding:"required"`
}

// ImportApps 批量导入应用到数据库
func ImportApps(c *gin.Context) {
	var req ImportAppsReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	if len(req.Apps) == 0 {
		response.Fail(c, "没有要导入的应用")
		return
	}

	now := time.Now().Unix()
	var maxWeigh int
	database.DB.Model(&model.Category{}).Select("COALESCE(MAX(weigh),0)").Scan(&maxWeigh)

	count := 0
	for i, app := range req.Apps {
		name := toString(app["name"])
		if name == "" {
			continue
		}
		item := model.Category{
			Name:       name,
			Nickname:   toString(app["version"]),
			Image:      toString(app["iconURL"]),
			Bt1a:       toString(app["downloadURL"]),
			Bt2a:       toString(app["size"]),
			Bt1b:       toString(app["tintColor"]),
			Bt2b:       toString(app["lock"]),
			Keywords:   toString(app["versionDescription"]),
			Flag:       toString(app["isLanZouCloud"]),
			Type:       "default",
			Status:     "normal",
			Createtime: now,
			Updatetime: now,
			Weigh:      maxWeigh + len(req.Apps) - i,
		}
		database.DB.Create(&item)
		count++
	}

	response.OKMsg(c, fmt.Sprintf("成功导入 %d 个应用", count), nil)
}

func toString(v interface{}) string {
	if v == nil {
		return ""
	}
	switch val := v.(type) {
	case string:
		return val
	case float64:
		return fmt.Sprintf("%g", val)
	default:
		return fmt.Sprintf("%v", val)
	}
}
