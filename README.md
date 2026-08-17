# SklandCheckin

森空岛（Skland）自动签到 / 检票脚本，支持多款游戏的一站式每日任务。

## 功能特性

- **多游戏签到**：明日方舟、终末地游戏内每日签到
- **森空岛检票**：6 个板块一键打卡，累计积分
- **成就数据**：泡姆泡姆、Ex Astris 游戏成就查询
- **玩家信息**：展示绑定角色的等级、主线进度等详情
- **三种登录方式**：密码 / 短信验证码 / 授权码登录
- **Token 持久化**：首次登录后自动保存，后续无需重复登录
- **日志记录**：同时输出到控制台和文件，按天滚动，自动清理
- **CI/CD 支持**：原生适配 GitHub Actions，自动检测 CI 环境并切换非交互模式

## 快速开始

### 环境要求

- Go 1.24+
- 可访问 `zonai.skland.com` 和 `web-api.skland.com` 的网络环境

### 构建

```bash
go build -o skland-checkin ./cmd/skland-checkin
```

或直接运行：

```bash
go run ./cmd/skland-checkin
```

### 使用

启动程序后按提示选择登录方式：

```
请选择登录方式:
1. 密码登录 (可能触发人机验证)
2. 手机验证码登录 (可能触发人机验证)
3. 授权码登录
0. 输入"0"退出
```

登录成功后会自动执行全部签到任务，完成后按回车退出。

> **授权码登录**：使用浏览器登录 [森空岛官网](https://www.skland.com) 后，访问 `https://web-api.skland.com/account/info/hg` 获取 JSON 内容，将整段内容粘贴到程序中即可。

## GitHub Actions 自动签到

项目内置 CI 模式：当检测到环境变量 `CI=true`（GitHub Actions 会自动设置）时，程序自动跳过交互式菜单，直接读取 `PHONE` 和 `PASSWORD` 环境变量执行密码登录。

### 配置步骤

1. **Fork 或克隆本仓库**到你的 GitHub 账号。

2. **配置 Secrets**：进入仓库的 `Settings → Secrets and variables → Actions → New repository secret`，添加以下两个密钥：

   | Secret 名称 | 说明 |
   |-------------|------|
   | `PHONE` | 森空岛账号绑定的手机号 |
   | `PASSWORD` | 森空岛账号密码 |

3. **启用 Actions**：在仓库的 `Actions` 页面确认 workflow 已启用。如果 Fork 自他人仓库，需要手动点击 `I understand my workflows, go ahead and enable them`。

4. **手动测试**：在 `Actions` 页面选择 `Skland每日自动签到` workflow，点击 `Run workflow` 手动触发一次，确认运行正常。

### 定时计划

Workflow 默认每天北京时间 00:00（UTC 16:00）自动执行，与森空岛签到/检票的 0 点刷新时间对齐。

> GitHub Actions 的定时任务可能有数分钟到数十分钟的延迟，属正常现象。如果需要更准点，可将 cron 改为提前几分钟触发。

### 日志查看

每次运行结束后，日志会作为 Artifact 上传，保留 3 天。在运行详情页的 `Artifacts` 区域下载 `logs-<run_id>` 即可查看。

### 注意事项

- **人机验证**：密码登录可能触发鹰角的人机验证，CI 环境下无法处理。如果频繁触发，建议改用本地运行 + 授权码登录，获取 token 后通过其他方式持久化。
- **Token 过期**：CI 模式下每次运行都会重新密码登录，不依赖缓存中的 token，因此无需关心 token 过期问题。
- **账号安全**：密码存储在 GitHub Secrets 中，仅在工作流运行时注入，不会出现在日志或代码中。

## 支持的游戏

### 游戏签到

| 游戏 | 说明 |
|------|------|
| 明日方舟 | 支持多绑定角色签到 |
| 明日方舟：终末地 | 按服务器角色分别签到 |

### 社区检票（积分任务）

| 板块 | 游戏 ID |
|------|---------|
| 明日方舟 | 1 |
| 来自星辰 | 2 |
| 明日方舟：终末地 | 3 |
| 泡姆泡姆 | 4 |
| 纳斯特港 | 100 |
| 开拓芯 | 101 |

### 成就查询（代码存在成就查询接口，但是没有使用）

| 游戏 | 说明 |
|------|------|
| 泡姆泡姆 | 查询游戏成就数据 |
| 来自星尘 | 查询剧情/结局完成情况 |

## 项目结构

```
.
├── cmd/skland-checkin/      # 程序入口
├── internal/
│   ├── api/                  # Skland API 封装（认证、签到、检票等）
│   ├── app/                  # 业务编排层（认证流程、签到流程）
│   ├── client/               # HTTP 客户端（签名、请求执行）
│   ├── config/               # 配置与 Token 持久化
│   ├── logger/               # 日志初始化（控制台 + 文件滚动）
│   ├── models/               # 数据模型定义
│   └── utils/                # 工具函数
├── .github/workflows/        # GitHub Actions CI 配置
└── go.mod
```

## License

MIT © 2025 HeathErika
