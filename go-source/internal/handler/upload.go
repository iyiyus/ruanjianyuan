package handler

import (
	"archive/zip"
	"bytes"
	"fmt"
	"go-source/internal/database"
	"go-source/internal/model"
	"go-source/internal/response"
	"go-source/internal/storage"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"howett.net/plist"
)

// UploadFile 上传图片（本地或云存储）
func UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, "请选择文件")
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".gif": true, ".webp": true}
	if !allowed[ext] {
		response.Fail(c, "只支持图片格式")
		return
	}

	src, err := file.Open()
	if err != nil {
		response.Fail(c, "读取文件失败")
		return
	}
	defer src.Close()
	data, _ := io.ReadAll(src)

	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	cfg := getStorageConfig()
	fileURL, err := storage.Upload(cfg, filename, data, file.Header.Get("Content-Type"))
	if err != nil {
		response.Fail(c, "上传失败: "+err.Error())
		return
	}
	response.OK(c, gin.H{"url": fileURL})
}

// UploadIPA 上传IPA并解析信息
func UploadIPA(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		response.Fail(c, "请选择IPA文件")
		return
	}
	if !strings.HasSuffix(strings.ToLower(file.Filename), ".ipa") {
		response.Fail(c, "只支持.ipa文件")
		return
	}

	src, err := file.Open()
	if err != nil {
		response.Fail(c, "读取文件失败")
		return
	}
	defer src.Close()
	data, _ := io.ReadAll(src)

	// 解析IPA（zip格式）
	info, iconData, iconExt, parseErr := parseIPA(data)

	// 上传IPA
	ipaFilename := fmt.Sprintf("%d.ipa", time.Now().UnixNano())
	cfg := getStorageConfig()
	ipaURL, uploadErr := storage.Upload(cfg, ipaFilename, data, "application/octet-stream")
	if uploadErr != nil {
		response.Fail(c, "IPA上传失败: "+uploadErr.Error())
		return
	}

	result := gin.H{
		"downloadURL": ipaURL,
		"size":        fmt.Sprintf("%.2f", float64(len(data))/1024/1024),
	}

	if parseErr == nil && info != nil {
		result["name"] = info["CFBundleDisplayName"]
		if result["name"] == "" {
			result["name"] = info["CFBundleName"]
		}
		result["version"] = info["CFBundleShortVersionString"]
		result["identifier"] = info["CFBundleIdentifier"]

		// 上传图标
		if iconData != nil {
			iconFilename := fmt.Sprintf("%d%s", time.Now().UnixNano(), iconExt)
			iconURL, err := storage.Upload(cfg, iconFilename, iconData, "image/png")
			if err == nil {
				result["iconURL"] = iconURL
			}
		}
	}

	response.OK(c, result)
}

// parseIPA 解析IPA文件，返回Info.plist内容和图标数据
func parseIPA(data []byte) (map[string]string, []byte, string, error) {
	r, err := zip.NewReader(bytes.NewReader(data), int64(len(data)))
	if err != nil {
		return nil, nil, "", err
	}

	var infoPlist *zip.File
	var iconFile *zip.File
	var iconExt string

	for _, f := range r.File {
		name := f.Name
		// 找 Info.plist：Payload/xxx.app/Info.plist
		if infoPlist == nil && strings.Count(name, "/") == 2 &&
			strings.HasSuffix(name, "/Info.plist") &&
			strings.HasPrefix(name, "Payload/") {
			infoPlist = f
		}
		// 找图标：优先 AppIcon60x60@2x.png
		if iconFile == nil && strings.HasPrefix(name, "Payload/") {
			base := filepath.Base(name)
			if strings.Contains(base, "AppIcon60x60@2x") && strings.HasSuffix(base, ".png") {
				iconFile = f
				iconExt = ".png"
			} else if iconFile == nil && strings.Contains(base, "AppIcon") && strings.HasSuffix(base, ".png") {
				iconFile = f
				iconExt = ".png"
			}
		}
	}

	if infoPlist == nil {
		return nil, nil, "", fmt.Errorf("未找到Info.plist")
	}

	// 读取 Info.plist
	rc, err := infoPlist.Open()
	if err != nil {
		return nil, nil, "", err
	}
	defer rc.Close()
	plistData, _ := io.ReadAll(rc)

	var rawInfo map[string]interface{}
	if _, err := plist.Unmarshal(plistData, &rawInfo); err != nil {
		return nil, nil, "", err
	}

	info := make(map[string]string)
	for k, v := range rawInfo {
		if s, ok := v.(string); ok {
			info[k] = s
		}
	}

	// 读取图标
	var iconData []byte
	if iconFile != nil {
		rc2, err := iconFile.Open()
		if err == nil {
			defer rc2.Close()
			iconData, _ = io.ReadAll(rc2)
		}
	}

	return info, iconData, iconExt, nil
}

// getStorageConfig 从数据库读取存储配置
func getStorageConfig() storage.Config {
	var list []model.Config
	database.DB.Find(&list)
	m := make(map[string]string)
	for _, v := range list {
		m[v.Name] = v.Value
	}
	return storage.GetConfigFromDB(m)
}
