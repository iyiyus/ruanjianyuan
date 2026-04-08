# iOS 软件源管理系统

基于 Go + Vue3 开发的 iOS 软件源管理后台，功能与 PHP 版本完全一致。

## 技术栈

- 后端：Go 1.21+ / Gin / GORM / MySQL
- 前端：Vue3 / TypeScript / Element Plus / Vite
- 部署：宝塔面板 / Linux

## 功能

- App 管理（支持 IPA 上传自动解析、云存储）
- 卡密管理（生成/激活/批量删除）
- UDID 黑名单（含公开对接接口）
- UDID 监控
- 系统配置（软件源/网站/云存储）
- 工作台数据统计
- 软件源搬运
- 安装引导系统

## 安装部署

### 环境要求

- Linux 服务器（推荐宝塔面板）
- Go 1.21+
- MySQL 5.7+
- Node.js 20+（仅编译前端需要）

---

### 方式一：使用编译好的二进制（推荐）

1. 下载 Release 中的 `ruanjianyuan.zip`，上传到服务器并解压

2. 进入解压目录，直接运行：
   ```bash
   ./go-source
   ```

3. 浏览器访问 `http://服务器IP:1117/install` 进入安装引导

4. 按引导填写数据库信息、管理员账号、后台路径等，完成安装

5. 安装完成后程序自动重启，访问 `http://服务器IP:端口/你设置的后台路径` 进入后台

---

### 方式二：源码编译

#### 编译前端

```bash
cd 中后台
pnpm install
pnpm build
cp -r dist ../go-source/dist
```

#### 编译后端

```bash
cd go-source
# Linux
GOOS=linux GOARCH=amd64 go build -o go-source-linux .
```

#### 运行

```bash
./go-source-linux
```

---

### 宝塔面板部署

1. 宝塔 → 软件商店 → 安装 **Go 项目管理器**

2. 新建 Go 项目：
   - 项目路径：上传的程序目录
   - 启动文件：`go-source-linux`
   - 端口：`1117`（或自定义）

3. 添加网站，开启反向代理：
   ```
   代理目标：http://127.0.0.1:1117
   ```

4. 访问域名 `/install` 完成安装引导

---

### 安装引导说明

访问 `/install` 路径，按步骤填写：

| 配置项 | 说明 |
|--------|------|
| 数据库主机 | 通常为 `127.0.0.1` |
| 数据库端口 | 默认 `3306` |
| 数据库用户名 | MySQL 用户名 |
| 数据库密码 | MySQL 密码 |
| 数据库名 | 提前创建好的数据库名 |
| 服务端口 | Go 程序监听端口，默认 `1117` |
| 后台路径 | 自定义后台入口，如 `admin` |
| GO项目名称 | 宝塔中的项目名，用于自动重启 |
| 网站名称 | 显示在后台的站点名称 |
| 管理员账号 | 登录后台的用户名 |
| 管理员密码 | 登录后台的密码 |

---

## 软件源接口

公开接口，无需认证：

```
GET /appstore        # 获取软件源数据
GET /Blacklist       # 获取黑名单列表
GET /Blacklist/check?udid=xxx  # 检查单个 UDID
```

## 目录结构

```
├── go-source/          # Go 后端源码
│   ├── main.go
│   ├── config/         # 配置读取
│   ├── router/         # 路由
│   ├── internal/
│   │   ├── handler/    # 接口处理
│   │   ├── model/      # 数据模型
│   │   ├── middleware/ # 中间件
│   │   └── storage/    # 云存储
│   └── dist/           # 前端构建产物（运行时需要）
└── 中后台/              # Vue3 前端源码
    ├── src/
    └── package.json
```

## License

MIT
