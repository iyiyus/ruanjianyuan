package handler

import (
	"go-source/internal/response"

	"github.com/gin-gonic/gin"
)

func GetMenus(c *gin.Context) {
	menus := []map[string]interface{}{
		{
			"path":      "/dashboard/console",
			"name":      "Console",
			"component": "/dashboard/console",
			"meta": map[string]interface{}{
				"title":    "menus.dashboard.console",
				"icon":     "ri:home-smile-2-line",
				"fixedTab": true,
				"roles":    []string{"R_SUPER", "R_ADMIN"},
			},
		},
		{
			"path":      "/source/apps",
			"name":      "Apps",
			"component": "/source/apps",
			"meta": map[string]interface{}{
				"title":     "App管理",
				"icon":      "ri:app-store-line",
				"keepAlive": true,
				"roles":     []string{"R_SUPER", "R_ADMIN"},
			},
		},
		{
			"path":      "/source/config",
			"name":      "SourceConfig",
			"component": "/source/config",
			"meta": map[string]interface{}{
				"title":     "源配置",
				"icon":      "ri:settings-3-line",
				"keepAlive": true,
				"roles":     []string{"R_SUPER"},
			},
		},
		{
			"path":      "/kami/list",
			"name":      "KamiList",
			"component": "/kami/list",
			"meta": map[string]interface{}{
				"title":     "卡密管理",
				"icon":      "ri:key-2-line",
				"keepAlive": true,
				"roles":     []string{"R_SUPER"},
			},
		},
		{
			"path":      "/security/black",
			"name":      "Black",
			"component": "/security/black",
			"meta": map[string]interface{}{
				"title":     "UDID黑名单",
				"icon":      "ri:forbid-line",
				"keepAlive": true,
				"roles":     []string{"R_SUPER"},
			},
		},
		{
			"path":      "/security/monitor",
			"name":      "Monitor",
			"component": "/security/monitor",
			"meta": map[string]interface{}{
				"title":     "UDID监控",
				"icon":      "ri:eye-line",
				"keepAlive": true,
				"roles":     []string{"R_SUPER"},
			},
		},
		{
			"path":      "/system/user",
			"name":      "User",
			"component": "/system/user",
			"meta": map[string]interface{}{
				"title":     "menus.system.user",
				"icon":      "ri:user-line",
				"keepAlive": true,
				"roles":     []string{"R_SUPER", "R_ADMIN"},
			},
		},
	}
	response.OK(c, menus)
}
