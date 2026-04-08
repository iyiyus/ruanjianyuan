package handler

import (
	"fmt"
	"go-source/internal/database"
	"go-source/internal/model"
	"go-source/internal/response"
	"time"

	"github.com/gin-gonic/gin"
)

type MonthStat struct {
	Month string `gorm:"column:month"`
	Count int    `gorm:"column:count"`
}

func GetDashboard(c *gin.Context) {
	var appCount, kamiTotal, kamiUsed, blackCount, monitorCount int64

	database.DB.Model(&model.Category{}).Where("status = 'normal'").Count(&appCount)
	database.DB.Model(&model.Kami{}).Count(&kamiTotal)
	database.DB.Model(&model.Kami{}).Where("jh = 1").Count(&kamiUsed)
	database.DB.Model(&model.Black{}).Count(&blackCount)
	database.DB.Model(&model.Monitor{}).Count(&monitorCount)

	var kamiMonthly []MonthStat
	database.DB.Raw(`
		SELECT DATE_FORMAT(FROM_UNIXTIME(usetime), '%m月') as month, COUNT(*) as count
		FROM fa_kami
		WHERE jh = 1 AND usetime > UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 9 MONTH))
		GROUP BY DATE_FORMAT(FROM_UNIXTIME(usetime), '%Y%m')
		ORDER BY DATE_FORMAT(FROM_UNIXTIME(usetime), '%Y%m') ASC
	`).Scan(&kamiMonthly)

	var monitorMonthly []MonthStat
	database.DB.Raw(`
		SELECT DATE_FORMAT(FROM_UNIXTIME(addtime), '%m月') as month, COUNT(*) as count
		FROM fa_monitor
		WHERE addtime > UNIX_TIMESTAMP(DATE_SUB(NOW(), INTERVAL 12 MONTH))
		GROUP BY DATE_FORMAT(FROM_UNIXTIME(addtime), '%Y%m')
		ORDER BY DATE_FORMAT(FROM_UNIXTIME(addtime), '%Y%m') ASC
	`).Scan(&monitorMonthly)

	kamiLabels, kamiData := fillMonthData(kamiMonthly, 9)
	monitorLabels, monitorData := fillMonthData(monitorMonthly, 12)

	response.OK(c, gin.H{
		"appCount":       appCount,
		"kamiTotal":      kamiTotal,
		"kamiUsed":       kamiUsed,
		"blackCount":     blackCount,
		"monitorCount":   monitorCount,
		"kamiLabels":     kamiLabels,
		"kamiMonthly":    kamiData,
		"monitorLabels":  monitorLabels,
		"monitorMonthly": monitorData,
	})
}

func fillMonthData(stats []MonthStat, n int) ([]string, []int) {
	now := time.Now()
	labels := make([]string, n)
	data := make([]int, n)
	monthMap := make(map[string]int)
	for _, s := range stats {
		monthMap[s.Month] = s.Count
	}
	for i := 0; i < n; i++ {
		t := now.AddDate(0, -(n - 1 - i), 0)
		label := fmt.Sprintf("%02d月", t.Month())
		labels[i] = label
		data[i] = monthMap[label]
	}
	return labels, data
}
