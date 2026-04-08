package handler

import (
	"crypto/md5"
	"database/sql"
	"fmt"
	"go-source/internal/response"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

const installLock = "install.lock"

func IsInstalled() bool {
	_, err := os.Stat(installLock)
	return err == nil
}

// InstallStatus 检查安装状态
func InstallStatus(c *gin.Context) {
	response.OK(c, gin.H{"installed": IsInstalled()})
}

// InstallCheckDB 测试数据库连接
func InstallCheckDB(c *gin.Context) {
	var req struct {
		Host     string `json:"host"`
		Port     string `json:"port"`
		User     string `json:"user"`
		Password string `json:"password"`
		DBName   string `json:"dbName"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}
	if req.Port == "" {
		req.Port = "3306"
	}
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		req.User, req.Password, req.Host, req.Port, req.DBName)
	db, err := sql.Open("mysql", dsn)
	if err != nil || db.Ping() != nil {
		response.Fail(c, "数据库连接失败")
		return
	}
	db.Close()
	response.OKMsg(c, "连接成功", nil)
}

// InstallRun 执行安装
func InstallRun(c *gin.Context) {
	if IsInstalled() {
		response.Fail(c, "系统已安装")
		return
	}

	var req struct {
		DBHost      string `json:"dbHost"`
		DBPort      string `json:"dbPort"`
		DBUser      string `json:"dbUser"`
		DBPassword  string `json:"dbPassword"`
		DBName      string `json:"dbName"`
		AppPort     string `json:"appPort"`
		AdminPath   string `json:"adminPath"`
		ProjectName string `json:"projectName"`
		SiteName    string `json:"siteName"`
		AdminUser   string `json:"adminUser"`
		AdminPwd    string `json:"adminPwd"`
		AdminEmail  string `json:"adminEmail"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, "参数错误")
		return
	}

	if req.DBHost == "" || req.DBUser == "" || req.DBName == "" {
		response.Fail(c, "请填写完整的数据库信息")
		return
	}
	if req.DBPort == "" {
		req.DBPort = "3306"
	}
	if req.AppPort == "" {
		req.AppPort = "1117"
	}
	if req.AdminPath == "" {
		req.AdminPath = "admin"
	}
	if req.ProjectName == "" {
		req.ProjectName = "go-source"
	}
	if req.SiteName == "" {
		req.SiteName = "iOS软件源管理系统"
	}
	if req.AdminUser == "" {
		req.AdminUser = "admin"
	}
	if req.AdminPwd == "" {
		req.AdminPwd = "123456"
	}
	req.AdminPath = strings.Trim(req.AdminPath, "/")

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		req.DBUser, req.DBPassword, req.DBHost, req.DBPort, req.DBName)
	db, err := sql.Open("mysql", dsn)
	if err != nil || db.Ping() != nil {
		response.Fail(c, "数据库连接失败")
		return
	}
	defer db.Close()

	tables := []string{
		`CREATE TABLE IF NOT EXISTS fa_admin (id int(10) unsigned NOT NULL AUTO_INCREMENT,username varchar(20) NOT NULL DEFAULT '',nickname varchar(50) NOT NULL DEFAULT '',password varchar(32) NOT NULL DEFAULT '',salt varchar(30) NOT NULL DEFAULT '',avatar varchar(255) NOT NULL DEFAULT '',email varchar(100) NOT NULL DEFAULT '',loginfailure tinyint(1) unsigned NOT NULL DEFAULT '0',logintime int(10) DEFAULT NULL,loginip varchar(50) DEFAULT NULL,createtime int(10) DEFAULT NULL,updatetime int(10) DEFAULT NULL,token varchar(59) NOT NULL DEFAULT '',status varchar(30) NOT NULL DEFAULT 'normal',PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS fa_category (id int(10) unsigned NOT NULL AUTO_INCREMENT,pid int(10) unsigned NOT NULL DEFAULT '0',type varchar(30) NOT NULL DEFAULT '',name varchar(30) NOT NULL DEFAULT '',nickname varchar(50) NOT NULL DEFAULT '',flag varchar(255) NOT NULL DEFAULT '',image varchar(100) NOT NULL DEFAULT '',keywords varchar(255) NOT NULL DEFAULT '',description varchar(255) NOT NULL DEFAULT '',diyname varchar(30) NOT NULL DEFAULT '',createtime int(10) DEFAULT NULL,updatetime int(10) DEFAULT NULL,weigh int(10) NOT NULL DEFAULT '0',status varchar(30) NOT NULL DEFAULT '',bt1a varchar(255) DEFAULT NULL,bt1b varchar(255) DEFAULT '018084',bt2a varchar(255) DEFAULT NULL,bt2b varchar(255) DEFAULT NULL,beizhu varchar(256) DEFAULT NULL,flag2 varchar(255) DEFAULT NULL,cs int(11) DEFAULT '0',cstime int(11) DEFAULT NULL,PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS fa_kami (id int(11) NOT NULL AUTO_INCREMENT,kami varchar(128) NOT NULL,udid varchar(128) DEFAULT NULL,jh int(1) NOT NULL DEFAULT '0',addtime int(11) NOT NULL DEFAULT '0',usetime int(11) NOT NULL DEFAULT '0',endtime int(11) NOT NULL DEFAULT '0',kmyp int(9) NOT NULL DEFAULT '1',custom_days int(11) NOT NULL DEFAULT '0',PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS fa_black (id int(11) NOT NULL AUTO_INCREMENT,udid varchar(128) DEFAULT NULL,reason varchar(255) NOT NULL DEFAULT '',addtime int(11) NOT NULL DEFAULT '0',usetime int(11) NOT NULL DEFAULT '0',endtime int(11) NOT NULL DEFAULT '0',PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		`CREATE TABLE IF NOT EXISTS fa_monitor (id int(11) NOT NULL AUTO_INCREMENT,udid varchar(255) NOT NULL,identity varchar(255) NOT NULL,count varchar(255) NOT NULL DEFAULT '0',addtime varchar(255) NOT NULL DEFAULT '0',PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4`,
		"CREATE TABLE IF NOT EXISTS fa_config (id int(10) unsigned NOT NULL AUTO_INCREMENT,name varchar(30) NOT NULL DEFAULT '',`group` varchar(30) NOT NULL DEFAULT '',title varchar(100) NOT NULL DEFAULT '',tip varchar(100) NOT NULL DEFAULT '',type varchar(30) NOT NULL DEFAULT '',value text NOT NULL,content text NOT NULL,rule varchar(100) NOT NULL DEFAULT '',extend varchar(255) NOT NULL DEFAULT '',PRIMARY KEY (id)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4",
	}
	for _, t := range tables {
		db.Exec(t)
	}

	now := time.Now().Unix()
	salt := "install"
	inner := fmt.Sprintf("%x", md5.Sum([]byte(req.AdminPwd)))
	pwdHash := fmt.Sprintf("%x", md5.Sum([]byte(inner+salt)))
	db.Exec(`INSERT IGNORE INTO fa_admin (username,nickname,password,salt,avatar,email,createtime,updatetime,status) VALUES (?,?,?,?,?,?,?,?,'normal')`,
		req.AdminUser, req.AdminUser, pwdHash, salt, "/uploads/avatar.png", req.AdminEmail, now, now)

	configs := [][]string{
		{"name", "basic", "站点名称", "string", req.SiteName},
		{"message", "basic", "公告内容", "text", "欢迎使用本软件源"},
		{"identifier", "basic", "源识别标符", "string", ""},
		{"sourceURL", "basic", "软件来源", "string", ""},
		{"sourceicon", "basic", "源图标", "string", ""},
		{"payURL", "basic", "解锁发卡地址", "string", ""},
		{"unlockURL", "basic", "解锁接口地址", "string", ""},
		{"opencry", "basic", "软件源加密", "switch", "0"},
		{"openblack", "basic", "自动拉黑添加者", "switch", "0"},
		{"openblack2", "basic", "自动拉黑破解者", "switch", "0"},
		{"site_name", "basic", "网站名称", "string", req.SiteName},
		{"site_logo", "basic", "网站LOGO", "string", ""},
		{"go_project_name", "basic", "GO项目名称", "string", req.ProjectName},
		{"storage_driver", "storage", "存储驱动", "string", "local"},
		{"admin_path", "basic", "后台路径", "string", req.AdminPath},
	}
	for _, cfg := range configs {
		db.Exec("INSERT IGNORE INTO fa_config (name,`group`,title,type,value,content,rule,extend) VALUES (?,?,?,?,?,'','','')",
			cfg[0], cfg[1], cfg[2], cfg[3], cfg[4])
	}

	envContent := fmt.Sprintf("DB_HOST=%s\nDB_PORT=%s\nDB_USER=%s\nDB_PASSWORD=%s\nDB_NAME=%s\nAPP_PORT=%s\nADMIN_PATH=%s\n",
		req.DBHost, req.DBPort, req.DBUser, req.DBPassword, req.DBName, req.AppPort, req.AdminPath)
	os.WriteFile(".env", []byte(envContent), 0644)
	os.WriteFile(installLock, []byte("installed at "+time.Now().Format("2006-01-02 15:04:05")), 0644)

	// 安装完成后自动重启（宝塔 Go 项目管理）
	go func() {
		time.Sleep(2 * time.Second)
		projectName := req.ProjectName
		if projectName == "" {
			projectName = "go-source"
		}
		// 尝试多种重启方式
		exec.Command("bt", "restart", projectName).Run()
		exec.Command("supervisorctl", "restart", projectName).Run()
	}()

	response.OK(c, gin.H{
		"adminUser": req.AdminUser,
		"adminPwd":  req.AdminPwd,
		"adminPath": req.AdminPath,
		"appPort":   req.AppPort,
		"siteName":  req.SiteName,
	})
}
