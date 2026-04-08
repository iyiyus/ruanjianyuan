package storage

import (
	"bytes"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type Config struct {
	Driver string // local | qiniu | tencent | aliyun | webdav
	// 七牛
	QiniuAccessKey string
	QiniuSecretKey string
	QiniuBucket    string
	QiniuDomain    string
	QiniuZone      string // z0 z1 z2 na0 as0
	// 腾讯云
	TencentSecretID  string
	TencentSecretKey string
	TencentBucket    string
	TencentRegion    string
	TencentDomain    string
	// 阿里云
	AliyunAccessKey string
	AliyunSecretKey string
	AliyunBucket    string
	AliyunEndpoint  string
	AliyunDomain    string
	// WebDAV
	WebdavURL      string
	WebdavUsername string
	WebdavPassword string
	WebdavDomain   string
}

// Upload 上传文件，返回可访问URL
func Upload(cfg Config, filename string, data []byte, contentType string) (string, error) {
	switch cfg.Driver {
	case "qiniu":
		return uploadQiniu(cfg, filename, data, contentType)
	case "tencent":
		return uploadTencent(cfg, filename, data, contentType)
	case "aliyun":
		return uploadAliyun(cfg, filename, data, contentType)
	case "webdav":
		return uploadWebdav(cfg, filename, data, contentType)
	default:
		return uploadLocal(filename, data)
	}
}

// ===== 本地存储 =====
func uploadLocal(filename string, data []byte) (string, error) {
	dir := "uploads"
	os.MkdirAll(dir, 0755)
	dst := filepath.Join(dir, filename)
	if err := os.WriteFile(dst, data, 0644); err != nil {
		return "", err
	}
	return "/uploads/" + filename, nil
}

// ===== 七牛云 =====
func uploadQiniu(cfg Config, filename string, data []byte, contentType string) (string, error) {
	// 生成上传token
	scope := cfg.QiniuBucket + ":" + filename
	deadline := time.Now().Unix() + 3600
	putPolicy := fmt.Sprintf(`{"scope":"%s","deadline":%d}`, scope, deadline)
	encoded := base64.URLEncoding.EncodeToString([]byte(putPolicy))
	mac := hmac.New(sha1.New, []byte(cfg.QiniuSecretKey))
	mac.Write([]byte(encoded))
	sign := base64.URLEncoding.EncodeToString(mac.Sum(nil))
	token := cfg.QiniuAccessKey + ":" + sign + ":" + encoded

	// 上传
	uploadURL := qiniuUploadURL(cfg.QiniuZone)
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("token", token)
	writer.WriteField("key", filename)
	part, _ := writer.CreateFormFile("file", filename)
	part.Write(data)
	writer.Close()

	resp, err := http.Post(uploadURL, writer.FormDataContentType(), body)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var result map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&result)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("七牛上传失败: %v", result)
	}
	domain := strings.TrimRight(cfg.QiniuDomain, "/")
	return domain + "/" + filename, nil
}

func qiniuUploadURL(zone string) string {
	zones := map[string]string{
		"z0":  "https://up.qiniup.com",
		"z1":  "https://up-z1.qiniup.com",
		"z2":  "https://up-z2.qiniup.com",
		"na0": "https://up-na0.qiniup.com",
		"as0": "https://up-as0.qiniup.com",
	}
	if u, ok := zones[zone]; ok {
		return u
	}
	return "https://up.qiniup.com"
}

// ===== 腾讯云 COS =====
func uploadTencent(cfg Config, filename string, data []byte, contentType string) (string, error) {
	host := fmt.Sprintf("%s.cos.%s.myqcloud.com", cfg.TencentBucket, cfg.TencentRegion)
	urlStr := "https://" + host + "/" + filename

	req, _ := http.NewRequest("PUT", urlStr, bytes.NewReader(data))
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Host", host)

	now := time.Now().UTC()
	date := now.Format("20060102")
	datetime := now.Format("20060102T150405Z")
	req.Header.Set("x-cos-date", datetime)

	// 签名
	signedHeaders := "content-type;host;x-cos-date"
	canonicalHeaders := fmt.Sprintf("content-type:%s\nhost:%s\nx-cos-date:%s\n", contentType, host, datetime)
	hashedPayload := sha256Hex(data)
	canonicalRequest := strings.Join([]string{"PUT", "/" + filename, "", canonicalHeaders, signedHeaders, hashedPayload}, "\n")
	credentialScope := date + "/cos/tc3_request"
	stringToSign := "TC3-HMAC-SHA256\n" + datetime + "\n" + credentialScope + "\n" + sha256Hex([]byte(canonicalRequest))

	secretDate := hmacSHA256([]byte("TC3"+cfg.TencentSecretKey), date)
	secretService := hmacSHA256(secretDate, "cos")
	secretSigning := hmacSHA256(secretService, "tc3_request")
	signature := hex.EncodeToString(hmacSHA256(secretSigning, stringToSign))

	auth := fmt.Sprintf("TC3-HMAC-SHA256 Credential=%s/%s, SignedHeaders=%s, Signature=%s",
		cfg.TencentSecretID, credentialScope, signedHeaders, signature)
	req.Header.Set("Authorization", auth)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("腾讯云上传失败: %s", string(b))
	}
	domain := cfg.TencentDomain
	if domain == "" {
		domain = "https://" + host
	}
	return strings.TrimRight(domain, "/") + "/" + filename, nil
}

// ===== 阿里云 OSS =====
func uploadAliyun(cfg Config, filename string, data []byte, contentType string) (string, error) {
	endpoint := cfg.AliyunEndpoint
	if !strings.HasPrefix(endpoint, "http") {
		endpoint = "https://" + endpoint
	}
	host := cfg.AliyunBucket + "." + strings.TrimPrefix(strings.TrimPrefix(endpoint, "https://"), "http://")
	urlStr := "https://" + host + "/" + filename

	now := time.Now().UTC()
	date := now.Format("Mon, 02 Jan 2006 15:04:05 GMT")
	contentMD5 := base64.StdEncoding.EncodeToString(md5Sum(data))

	stringToSign := strings.Join([]string{"PUT", contentMD5, contentType, date,
		"/" + cfg.AliyunBucket + "/" + filename}, "\n")
	mac := hmac.New(sha1.New, []byte(cfg.AliyunSecretKey))
	mac.Write([]byte(stringToSign))
	signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	req, _ := http.NewRequest("PUT", urlStr, bytes.NewReader(data))
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Content-MD5", contentMD5)
	req.Header.Set("Date", date)
	req.Header.Set("Authorization", "OSS "+cfg.AliyunAccessKey+":"+signature)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("阿里云上传失败: %s", string(b))
	}
	domain := cfg.AliyunDomain
	if domain == "" {
		domain = "https://" + host
	}
	return strings.TrimRight(domain, "/") + "/" + filename, nil
}

// ===== WebDAV =====
func uploadWebdav(cfg Config, filename string, data []byte, contentType string) (string, error) {
	urlStr := strings.TrimRight(cfg.WebdavURL, "/") + "/" + filename
	req, _ := http.NewRequest("PUT", urlStr, bytes.NewReader(data))
	req.Header.Set("Content-Type", contentType)
	if cfg.WebdavUsername != "" {
		req.SetBasicAuth(cfg.WebdavUsername, cfg.WebdavPassword)
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return "", fmt.Errorf("WebDAV上传失败: %d", resp.StatusCode)
	}
	domain := cfg.WebdavDomain
	if domain == "" {
		domain = cfg.WebdavURL
	}
	return strings.TrimRight(domain, "/") + "/" + filename, nil
}

// ===== 工具函数 =====
func sha256Hex(data []byte) string {
	h := sha256.Sum256(data)
	return hex.EncodeToString(h[:])
}

func hmacSHA256(key []byte, data string) []byte {
	mac := hmac.New(sha256.New, key)
	mac.Write([]byte(data))
	return mac.Sum(nil)
}

func md5Sum(data []byte) []byte {
	h := md5.Sum(data)
	return h[:]
}

// GetConfigFromDB 从配置map构建存储配置
func GetConfigFromDB(m map[string]string) Config {
	return Config{
		Driver:           m["storage_driver"],
		QiniuAccessKey:   m["qiniu_access_key"],
		QiniuSecretKey:   m["qiniu_secret_key"],
		QiniuBucket:      m["qiniu_bucket"],
		QiniuDomain:      m["qiniu_domain"],
		QiniuZone:        m["qiniu_zone"],
		TencentSecretID:  m["tencent_secret_id"],
		TencentSecretKey: m["tencent_secret_key"],
		TencentBucket:    m["tencent_bucket"],
		TencentRegion:    m["tencent_region"],
		TencentDomain:    m["tencent_domain"],
		AliyunAccessKey:  m["aliyun_access_key"],
		AliyunSecretKey:  m["aliyun_secret_key"],
		AliyunBucket:     m["aliyun_bucket"],
		AliyunEndpoint:   m["aliyun_endpoint"],
		AliyunDomain:     m["aliyun_domain"],
		WebdavURL:        m["webdav_url"],
		WebdavUsername:   m["webdav_username"],
		WebdavPassword:   m["webdav_password"],
		WebdavDomain:     m["webdav_domain"],
	}
}

// 避免unused import
var _ = sort.Strings
var _ = url.QueryEscape
