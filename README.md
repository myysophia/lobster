# Lobster

Lobster 是一个面向多种 Claw 系产品的统一安装入口项目。

当前第一优先级目标是：

- 用一个轻量的跨平台 CLI 大幅简化腾讯 WorkBuddy 的安装
- 让用户只需要记住一个安装命令
- 在安装完成后给出清晰的下一步引导，而不是让用户回到官网自行摸索

## 当前范围

第一版聚焦 WorkBuddy：

- 统一入口安装
- 安装结果校验
- 首次启动引导
- 常见错误诊断

后续预留支持：

- ArkClaw
- Kimi Claw
- AutoClaw
- 其他同类产品

## 当前建议

- 项目工作区名称使用 `lobster`
- 第一版实现语言使用 `Go`
- 第一版采用“薄封装”策略：调用官方安装器，不自建版本分发体系

## 目录结构

```text
lobster/
  README.md
  docs/
    plans/
      2026-03-22-lobster-design.md
      2026-03-22-lobster-implementation-plan.md
      2026-03-22-lobster-milestones.md
```

## 当前输出物

- 设计文档：`docs/plans/2026-03-22-lobster-design.md`
- 实施计划：`docs/plans/2026-03-22-lobster-implementation-plan.md`
- 里程碑：`docs/plans/2026-03-22-lobster-milestones.md`

## 当前 CLI 约定

主命令：

```bash
lobster tui
lobster install workbuddy
lobster status workbuddy
lobster open workbuddy
lobster doctor workbuddy
lobster list
```

快捷别名：

```bash
wb tui
wb install
wb status
wb open
wb doctor
wb next
```

## 当前实现状态

截至 2026-03-23，当前原型已经具备以下行为：

- `lobster install workbuddy --dry-run` 会输出平台识别结果、安装策略与官方安装命令，但不会真正执行安装
- `lobster install workbuddy` 会先做安装前检测；如果当前已经检测到可用命令，会跳过重复安装
- 安装命令执行后会立即做安装后复检，并根据结果输出下一步建议
- `lobster status workbuddy` 会区分“已检测到可用安装”“已检测到安装痕迹但命令暂不可用”“未检测到安装”
- `lobster doctor workbuddy` 会输出命令可用性、命中路径、环境提示与建议操作
- `wb next`、`wb open` 等快捷别名与主命令保持并行可用

当前已补充的基础验证包括：

```bash
go test ./...
go build ./...
```

## TUI 安装向导

项目现已支持一个面向新手的终端安装向导。

入口命令：

```bash
lobster tui
wb tui
```

当前行为：

- `lobster tui` 会进入多产品选择页
- `wb tui` 会直接进入 `WorkBuddy` 安装向导
- `WorkBuddy` 已接入真实安装流程
- `ArkClaw`、`Kimi Claw`、`AutoClaw` 暂只显示 `On The Way`

在多产品页中，用户可以用方向键或 `j/k` 切换目标，按下 `Enter` 进入 WorkBuddy 向导或跳转到预留的 `On The Way` 占位页，而未完成的产品始终不会触发真正的安装逻辑，只给出即将开放的提示。`wb tui` 则跳过列表，直接定位 WorkBuddy 安装体验。

WorkBuddy 向导内部会先执行状态探测，再通过 `installer.RunWithIO` 直接运行官方安装器，贴合原生输出；安装中页会在完成后自动转到结果页面，结果页不仅展示依赖 `installer.Result` 的 Outcome，还呈现最新的安装输出、下一步建议，以及可选的诊断详情。

当前 TUI 支持的交互：

- 产品列表中使用 `↑/↓` 或 `j/k` 切换
- `Enter` 进入产品或执行安装
- `Esc` 返回上一级
- `r` 重新检查状态
- `o` 打开应用
- `d` 查看诊断详情
- `q` 退出

适用场景：

- 首次安装和体验 `WorkBuddy`
- 需要更强引导感的终端交互

不建议在以下场景使用 TUI：

- CI
- shell 脚本
- 非交互式终端

这些场景仍建议直接使用现有 CLI 子命令。

## GitHub Release

仓库已支持通过 GitHub Actions 在打 tag 时自动发布多平台构建产物。

触发方式：

```bash
git tag v0.1.0
git push origin v0.1.0
```

当前发布平台：

- `darwin/amd64`
- `darwin/arm64`
- `linux/amd64`
- `linux/arm64`
- `windows/amd64`

发布产物规则：

- 包名格式：`lobster_<version>_<os>_<arch>.tar.gz`
- Windows 包名格式：`lobster_<version>_<os>_<arch>.zip`
- 每个压缩包内同时包含 `lobster`、`wb` 和 `README.md`
- Release 中会同时上传对应的 `*.sha256` 文件与汇总 `SHA256SUMS`

本地手动打包：

```bash
GOOS=darwin GOARCH=arm64 VERSION=v0.1.0 ./scripts/build-release.sh
GOOS=linux GOARCH=amd64 VERSION=v0.1.0 ./scripts/build-release.sh
GOOS=windows GOARCH=amd64 VERSION=v0.1.0 ./scripts/build-release.sh
```

本地产物默认输出到：

```bash
dist/
```

## 本地开发

运行主命令：

```bash
go run ./cmd/lobster help
go run ./cmd/lobster tui
go run ./cmd/lobster list
go run ./cmd/lobster install workbuddy --dry-run
```

运行 WorkBuddy 快捷别名：

```bash
go run ./cmd/wb help
go run ./cmd/wb tui
go run ./cmd/wb install --dry-run
```

编译：

```bash
go build ./...
```
