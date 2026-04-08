package handler

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"go-source/internal/database"
	"go-source/internal/model"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// AppStore 对外iOS软件源接口，完全兼容PHP版本
func AppStore(c *gin.Context) {
	nowtime := time.Now().Format("2006-01-02 15:04:05")

	// 获取当前域名，用于补全相对路径
	scheme := "https"
	if c.Request.TLS == nil {
		scheme = "http"
	}
	baseURL := scheme + "://" + c.Request.Host

	// 读取POST body（黑名单上报）
	body, _ := io.ReadAll(c.Request.Body)
	if len(body) > 0 {
		var postData map[string]interface{}
		if json.Unmarshal(body, &postData) == nil {
			if val, ok := postData["value"].(string); ok && val != "" {
				decoded, err := base64.StdEncoding.DecodeString(val)
				if err == nil {
					parts := strings.Split(string(decoded), "|")
					if len(parts) >= 2 {
						udid1 := parts[0]
						udid2 := parts[1]
						handleUdidReport(udid1, udid2)
					}
				}
			}
		}
	}

	udid := c.Query("udid")
	kcode := c.Query("code")

	// 检查黑名单
	if udid != "" {
		var black model.Black
		if database.DB.Where("udid = ?", udid).First(&black).Error == nil {
			result := buildBlackResponse(udid, nowtime)
			c.JSON(http.StatusOK, result)
			return
		}
	}

	// 获取配置
	cfgMap := getConfigMap()
	opencry := cfgMap["opencry"]

	if kcode == "" {
		appList := buildAppList(udid, cfgMap, baseURL)
		info := buildSourceInfo(cfgMap)
		arr := buildResponse(info, appList, udid, nowtime)

		if opencry == "1" {
			b, _ := json.Marshal(arr)
			encoded := base64.StdEncoding.EncodeToString(b)
			c.JSON(http.StatusOK, gin.H{"appstore": encoded})
		} else {
			c.JSON(http.StatusOK, arr)
		}
		return
	}

	// 有卡密，验证并激活
	var kami model.Kami
	if err := database.DB.Where("kami = ?", kcode).Order("id DESC").First(&kami).Error; err != nil {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "解锁码不存在"})
		return
	}
	if kami.Jh == 1 {
		c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "解锁码已使用"})
		return
	}

	now := time.Now().Unix()
	var endtm int64
	switch kami.Kmyp {
	case 1:
		endtm = now + 86400*30
	case 2:
		endtm = now + 86400*90
	case 3:
		endtm = now + 86400*365
	case 4:
		days := kami.CustomDays
		if days <= 0 {
			days = 30
		}
		endtm = now + int64(86400*days)
	}
	database.DB.Model(&kami).Updates(map[string]interface{}{
		"udid":    udid,
		"usetime": now,
		"endtime": endtm,
		"jh":      1,
	})
	c.JSON(http.StatusOK, gin.H{"code": 0, "msg": "ok，解锁成功"})
}

func handleUdidReport(udid1, udid2 string) {
	cfgMap := getConfigMap()
	openblack := cfgMap["openblack"]
	openblack2 := cfgMap["openblack2"]

	processUdid := func(udid, identity string, autoBlack bool) {
		if udid == "" {
			return
		}
		l := len(udid)
		if l != 25 && l != 40 {
			return
		}
		if autoBlack {
			var exist model.Black
			if database.DB.Where("udid = ?", udid).First(&exist).Error != nil {
				database.DB.Create(&model.Black{Udid: udid, Addtime: time.Now().Unix()})
			}
			database.DB.Where("udid = ?", udid).Delete(&model.Monitor{})
		} else {
			var existBlack model.Black
			if database.DB.Where("udid = ?", udid).First(&existBlack).Error != nil {
				var m model.Monitor
				if database.DB.Where("udid = ?", udid).First(&m).Error == nil {
					database.DB.Model(&m).UpdateColumn("count", database.DB.Raw("count + 1"))
				} else {
					database.DB.Create(&model.Monitor{
						Udid:     udid,
						Identity: identity,
						Count:    "1",
						Addtime:  strings.Split(time.Now().Format("2006-01-02 15:04:05"), " ")[0],
					})
				}
			}
		}
	}

	processUdid(udid1, "添加者", openblack == "1")
	processUdid(udid2, "破解者", openblack2 == "1")
}

func buildAppList(udid string, cfgMap map[string]string, baseURL string) []map[string]interface{} {
	var categories []model.Category
	database.DB.Where("status = 'normal'").Order("weigh DESC").Find(&categories)

	hasValidKami := false
	if udid != "" {
		var kami model.Kami
		if database.DB.Where("udid = ? AND jh = 1 AND endtime > ?", udid, time.Now().Unix()).First(&kami).Error == nil {
			hasValidKami = true
		}
	}

	fullURL := func(u string) string {
		if u != "" && strings.HasPrefix(u, "/") {
			return baseURL + u
		}
		return u
	}

	var data []map[string]interface{}
	for _, v := range categories {
		t := 0
		if v.Type != "default" && v.Type != "" {
			switch v.Type {
			case "1":
				t = 1
			case "2":
				t = 2
			case "3":
				t = 3
			case "4":
				t = 4
			case "5":
				t = 5
			}
		}
		downloadURL := fullURL(v.Bt1a)
		if v.Bt2b == "1" && !hasValidKami {
			downloadURL = ""
		}
		data = append(data, map[string]interface{}{
			"name":               v.Name,
			"type":               t,
			"version":            v.Nickname,
			"versionDate":        time.Unix(v.Updatetime, 0).Format("2006-01-02T15:04:05+08:00"),
			"versionDescription": strings.ReplaceAll(v.Keywords, "\\n", "\n"),
			"lock":               v.Bt2b,
			"downloadURL":        downloadURL,
			"isLanZouCloud":      v.Flag,
			"iconURL":            fullURL(v.Image),
			"tintColor":          v.Bt1b,
			"size":               v.Bt2a,
		})
	}
	if data == nil {
		data = []map[string]interface{}{}
	}
	return data
}

func buildSourceInfo(cfgMap map[string]string) map[string]string {
	return map[string]string{
		"name":       cfgMap["name"],
		"message":    cfgMap["message"],
		"identifier": cfgMap["identifier"],
		"sourceURL":  cfgMap["sourceURL"],
		"sourceicon": cfgMap["sourceicon"],
		"payURL":     cfgMap["payURL"],
		"unlockURL":  cfgMap["unlockURL"],
	}
}

type AppItem struct {
	PayURL     string      `json:"payURL"`
	Apps       interface{} `json:"apps"`
	UnlockURL  string      `json:"unlockURL"`
	Message    string      `json:"message"`
	Identifier string      `json:"identifier"`
	SourceIcon string      `json:"sourceicon"`
	Name       string      `json:"name"`
	SourceURL  string      `json:"sourceURL"`
}

type AppEntry struct {
	IsLanZouCloud string `json:"isLanZouCloud"`
	VersionDate   string `json:"versionDate"`
	Lock          string `json:"lock"`
	IconURL       string `json:"iconURL"`
	VersionDesc   string `json:"versionDescription"`
	DownloadURL   string `json:"downloadURL"`
	Version       string `json:"version"`
	Type          int    `json:"type"`
	TintColor     string `json:"tintColor"`
	Name          string `json:"name"`
	Size          string `json:"size"`
}

func buildResponse(info map[string]string, apps []map[string]interface{}, udid, nowtime string) interface{} {
	entries := make([]AppEntry, 0, len(apps))
	for _, a := range apps {
		t := 0
		if v, ok := a["type"].(int); ok {
			t = v
		}
		entries = append(entries, AppEntry{
			IsLanZouCloud: str(a["isLanZouCloud"]),
			VersionDate:   str(a["versionDate"]),
			Lock:          str(a["lock"]),
			IconURL:       str(a["iconURL"]),
			VersionDesc:   str(a["versionDescription"]),
			DownloadURL:   str(a["downloadURL"]),
			Version:       str(a["version"]),
			Type:          t,
			TintColor:     str(a["tintColor"]),
			Name:          str(a["name"]),
			Size:          mbToBytes(str(a["size"])),
		})
	}
	return AppItem{
		PayURL:     info["payURL"],
		Apps:       entries,
		UnlockURL:  info["unlockURL"],
		Message:    info["message"],
		Identifier: info["identifier"],
		SourceIcon: info["sourceicon"],
		Name:       info["name"],
		SourceURL:  info["sourceURL"],
	}
}

func str(v interface{}) string {
	if v == nil {
		return ""
	}
	if s, ok := v.(string); ok {
		return s
	}
	return fmt.Sprintf("%v", v)
}

// mbToBytes 将MB字符串转为字节字符串（签名工具需要字节数）
func mbToBytes(mb string) string {
	if mb == "" {
		return "0"
	}
	var f float64
	fmt.Sscanf(mb, "%f", &f)
	bytes := int64(f * 1024 * 1024)
	return fmt.Sprintf("%d", bytes)
}

func buildBlackResponse(udid, nowtime string) map[string]interface{} {
	return map[string]interface{}{
		"name":       "已被源主拉黑",
		"message":    "你已被源主拉黑！",
		"identifier": "长按此处删除软件源",
		"payURL":     "",
		"unlockURL":  "",
		"UDID":       udid,
		"Time":       nowtime,
		"apps": []map[string]interface{}{
			{
				"name":               "你已被源主拉黑！",
				"version":            "1.0",
				"type":               "1.0",
				"versionDate":        "2021-01-24",
				"versionDescription": "你已被源主拉黑！",
				"lock":               "1",
				"downloadURL":        "",
				"isLanZouCloud":      "0",
				"tintColor":          "",
				"size":               "123973140.48",
			},
		},
	}
}

func sendEncrypted(c *gin.Context, data map[string]interface{}) {
	b, _ := json.Marshal(data)
	encoded := base64.StdEncoding.EncodeToString(b)
	c.JSON(http.StatusOK, gin.H{"appstore": encoded})
}

func getConfigMap() map[string]string {
	var list []model.Config
	database.DB.Find(&list)
	m := make(map[string]string)
	for _, v := range list {
		m[v.Name] = v.Value
	}
	return m
}
