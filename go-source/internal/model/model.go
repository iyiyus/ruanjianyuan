package model

// Admin 管理员表
type Admin struct {
	ID           uint   `gorm:"primaryKey;column:id"`
	Username     string `gorm:"column:username"`
	Nickname     string `gorm:"column:nickname"`
	Password     string `gorm:"column:password"`
	Salt         string `gorm:"column:salt"`
	Avatar       string `gorm:"column:avatar"`
	Email        string `gorm:"column:email"`
	Loginfailure int    `gorm:"column:loginfailure"`
	Logintime    int64  `gorm:"column:logintime"`
	Loginip      string `gorm:"column:loginip"`
	Createtime   int64  `gorm:"column:createtime"`
	Updatetime   int64  `gorm:"column:updatetime"`
	Token        string `gorm:"column:token"`
	Status       string `gorm:"column:status"`
}

func (Admin) TableName() string { return "fa_admin" }

// Category App分类表
type Category struct {
	ID          uint   `gorm:"primaryKey;column:id" json:"id"`
	Pid         uint   `gorm:"column:pid" json:"pid"`
	Type        string `gorm:"column:type" json:"type"`
	Name        string `gorm:"column:name" json:"name"`
	Nickname    string `gorm:"column:nickname" json:"nickname"`
	Flag        string `gorm:"column:flag" json:"flag"`
	Image       string `gorm:"column:image" json:"image"`
	Keywords    string `gorm:"column:keywords" json:"keywords"`
	Description string `gorm:"column:description" json:"description"`
	Diyname     string `gorm:"column:diyname" json:"diyname"`
	Createtime  int64  `gorm:"column:createtime" json:"createtime"`
	Updatetime  int64  `gorm:"column:updatetime" json:"updatetime"`
	Weigh       int    `gorm:"column:weigh" json:"weigh"`
	Status      string `gorm:"column:status" json:"status"`
	Bt1a        string `gorm:"column:bt1a" json:"bt1a"` // 下载链接
	Bt1b        string `gorm:"column:bt1b" json:"bt1b"` // 颜色
	Bt2a        string `gorm:"column:bt2a" json:"bt2a"` // 文件大小
	Bt2b        string `gorm:"column:bt2b" json:"bt2b"` // 是否加锁
	Beizhu      string `gorm:"column:beizhu" json:"beizhu"`
	Flag2       string `gorm:"column:flag2" json:"flag2"`
	Cs          int    `gorm:"column:cs" json:"cs"`
	Cstime      *int64 `gorm:"column:cstime" json:"cstime"`
}

func (Category) TableName() string { return "fa_category" }

// Kami 卡密表
type Kami struct {
	ID         int    `gorm:"primaryKey;column:id" json:"id"`
	Kami       string `gorm:"column:kami" json:"kami"`
	Udid       string `gorm:"column:udid" json:"udid"`
	Jh         int    `gorm:"column:jh" json:"jh"`
	Addtime    int64  `gorm:"column:addtime" json:"addtime"`
	Usetime    int64  `gorm:"column:usetime" json:"usetime"`
	Endtime    int64  `gorm:"column:endtime" json:"endtime"`
	Kmyp       int    `gorm:"column:kmyp" json:"kmyp"` // 1=30天 2=90天 3=365天 4=自定义
	CustomDays int    `gorm:"column:custom_days" json:"custom_days"`
}

func (Kami) TableName() string { return "fa_kami" }

// Black 黑名单表
type Black struct {
	ID      int    `gorm:"primaryKey;column:id" json:"id"`
	Udid    string `gorm:"column:udid" json:"udid"`
	Reason  string `gorm:"column:reason" json:"reason"`
	Addtime int64  `gorm:"column:addtime" json:"addtime"`
	Usetime int64  `gorm:"column:usetime" json:"usetime"`
	Endtime int64  `gorm:"column:endtime" json:"endtime"`
}

func (Black) TableName() string { return "fa_black" }

// Monitor UDID监控表
type Monitor struct {
	ID       int    `gorm:"primaryKey;column:id" json:"id"`
	Udid     string `gorm:"column:udid" json:"udid"`
	Identity string `gorm:"column:identity" json:"identity"`
	Count    string `gorm:"column:count" json:"count"`
	Addtime  string `gorm:"column:addtime" json:"addtime"`
}

func (Monitor) TableName() string { return "fa_monitor" }

// Config 系统配置表
type Config struct {
	ID      uint   `gorm:"primaryKey;column:id" json:"id"`
	Name    string `gorm:"column:name" json:"name"`
	Group   string `gorm:"column:group" json:"group"`
	Title   string `gorm:"column:title" json:"title"`
	Tip     string `gorm:"column:tip" json:"tip"`
	Type    string `gorm:"column:type" json:"type"`
	Value   string `gorm:"column:value" json:"value"`
	Content string `gorm:"column:content" json:"content"`
	Rule    string `gorm:"column:rule" json:"rule"`
	Extend  string `gorm:"column:extend" json:"extend"`
}

func (Config) TableName() string { return "fa_config" }

// AdminLog 管理员日志
type AdminLog struct {
	ID         uint   `gorm:"primaryKey;column:id" json:"id"`
	AdminID    uint   `gorm:"column:admin_id" json:"admin_id"`
	Username   string `gorm:"column:username" json:"username"`
	URL        string `gorm:"column:url" json:"url"`
	Title      string `gorm:"column:title" json:"title"`
	Content    string `gorm:"column:content" json:"content"`
	IP         string `gorm:"column:ip" json:"ip"`
	Useragent  string `gorm:"column:useragent" json:"useragent"`
	Createtime int64  `gorm:"column:createtime" json:"createtime"`
}

func (AdminLog) TableName() string { return "fa_admin_log" }
