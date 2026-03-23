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
lobster install workbuddy
lobster status workbuddy
lobster open workbuddy
lobster doctor workbuddy
lobster list
```

快捷别名：

```bash
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
go run ./cmd/lobster list
go run ./cmd/lobster install workbuddy --dry-run
```

运行 WorkBuddy 快捷别名：

```bash
go run ./cmd/wb help
go run ./cmd/wb install --dry-run
```

编译：

```bash
go build ./...
```
