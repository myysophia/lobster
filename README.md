# Lobster

Lobster 是一个面向多种 Claw 系产品的统一安装入口项目。

当前第一优先级目标是：

- 用一个轻量的跨平台 CLI 大幅简化腾讯 WorkBuddy 的安装
- 让用户只需要记住一个安装命令
- 在安装完成后给出清晰的下一步引导，而不是让用户回到官网自行摸索

## 当前范围

当前已落地能力：

- WorkBuddy：统一入口安装、安装结果校验、首次启动引导、常见错误诊断
- AutoClaw：按平台拉起官方安装包下载链接，覆盖 Windows / macOS Apple Silicon / macOS Intel
- QoderWork：按平台拉起官方安装包下载链接，覆盖 Windows / macOS Apple Silicon / macOS Intel

后续仍预留支持：

- ArkClaw
- Kimi Claw
- AutoClaw
- QoderWork
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

全局命令：

```bash
lobster help
lobster tui
lobster list
```

产品命令：

```bash
lobster workbuddy help
lobster workbuddy install [--dry-run]
lobster workbuddy status
lobster workbuddy open
lobster workbuddy doctor
lobster workbuddy next
lobster workbuddy tui
lobster autoclaw help
lobster autoclaw install [--dry-run]
lobster autoclaw status
lobster autoclaw open
lobster autoclaw doctor
lobster autoclaw next
lobster autoclaw tui
lobster qoderwork help
lobster qoderwork install [--dry-run]
lobster qoderwork status
lobster qoderwork open
lobster qoderwork doctor
lobster qoderwork next
lobster qoderwork tui
```

## 当前实现状态

截至 2026-04-02，当前原型已经具备以下行为：

- `lobster workbuddy install --dry-run` 会输出平台识别结果、安装策略与官方安装命令，但不会真正执行安装
- `lobster workbuddy install` 会先做安装前检测；如果当前已经检测到可用命令，会跳过重复安装
- 安装命令执行后会立即做安装后复检，并根据结果输出下一步建议
- `lobster workbuddy status` 会区分“已检测到可用安装”“已检测到安装痕迹但命令暂不可用”“未检测到安装”
- `lobster workbuddy doctor` 会输出命令可用性、命中路径、环境提示与建议操作
- `lobster autoclaw install --dry-run` 会按当前平台输出 AutoClaw 官方安装包直链；当前支持 `windows/amd64`、`darwin/arm64`、`darwin/amd64`
- `lobster qoderwork install --dry-run` 会按当前平台输出 QoderWork 官方安装包直链；当前支持 `windows/amd64`、`darwin/arm64`、`darwin/amd64`
- `lobster autoclaw install` 与 `lobster qoderwork install` 会直接拉起官方下载流程，并提示用户安装完成后执行 `status` 或 `doctor` 复查
- 默认安装输出已收敛为简洁模式，不再回显官方安装器全过程日志

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
lobster workbuddy tui
```

当前行为：

- `lobster tui` 会进入多产品选择页
- `lobster workbuddy tui` 会直接进入 `WorkBuddy` 安装向导
- `WorkBuddy` 已接入真实安装流程
- `ArkClaw`、`Kimi Claw` 暂只显示 `On The Way`
- `AutoClaw`、`QoderWork` 当前已在 CLI 中接入官方下载流程，但 TUI 里仍先保留 `On The Way`

在多产品页中，用户可以用方向键或 `j/k` 切换目标，按下 `Enter` 进入 WorkBuddy 向导或跳转到预留的 `On The Way` 占位页，而未完成的产品始终不会触发真正的安装逻辑，只给出即将开放的提示。`lobster workbuddy tui` 则跳过列表，直接定位 WorkBuddy 安装体验。

WorkBuddy 向导内部会先执行状态探测，再通过 `installer.RunWithIO` 直接运行官方安装器。安装中页保持简短提示，结果页默认展示简版结果与下一步建议，仅在失败或校验未通过时展示最近一次安装输出摘要。

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
- 每个压缩包内包含 `lobster` 和 `README.md`
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

## 按平台安装与使用

下面以“使用 Lobster 安装 WorkBuddy”为例，分别说明 Windows、macOS、Linux 的推荐使用方式。

通用原则：

- 先从 GitHub Release 下载与你当前系统和 CPU 架构匹配的压缩包
- 解压后进入目录，直接运行 `lobster` 即可
- 如果只是想先看命令会做什么，可以先执行 `--dry-run`
- 安装完成后建议继续执行 `status`、`next`、`open`、`doctor` 做确认

### Windows

适用产物示例：

- `lobster_<version>_windows_amd64.zip`

建议步骤（PowerShell）：

```powershell
# 1. 解压下载好的压缩包
Expand-Archive .\lobster_v0.1.0_windows_amd64.zip -DestinationPath .\lobster

# 2. 进入解压目录
cd .\lobster\lobster_v0.1.0_windows_amd64

# 3. 先查看帮助
.\lobster.exe help

# 4. 先做一次演练，不真正执行安装
.\lobster.exe workbuddy install --dry-run

# 5. 真正安装 WorkBuddy
.\lobster.exe workbuddy install

# 6. 安装完成后检查状态
.\lobster.exe workbuddy status

# 7. 查看下一步建议
.\\lobster.exe workbuddy next
```

Windows 使用提示：

- 如果你希望在任意目录直接输入 `lobster`，可以把解压目录加入 `PATH`
- 如果安装后命令暂时不可用，先执行 `.\lobster.exe workbuddy doctor` 查看诊断结果

### macOS

适用产物示例：

- `lobster_<version>_darwin_arm64.tar.gz`，适用于 Apple Silicon Mac
- `lobster_<version>_darwin_amd64.tar.gz`，适用于 Intel Mac

建议步骤（Terminal）：

```bash
# 1. 解压下载好的压缩包
tar -xzf lobster_v0.1.0_darwin_arm64.tar.gz

# 2. 进入解压目录
cd lobster_v0.1.0_darwin_arm64

# 3. 如有需要，补可执行权限
chmod +x lobster

# 4. 先查看帮助
./lobster help

# 5. 先做一次演练，不真正执行安装
./lobster workbuddy install --dry-run

# 6. 真正安装 WorkBuddy
./lobster workbuddy install

# 7. 安装完成后检查状态
./lobster workbuddy status

# 8. 查看下一步建议或尝试打开应用
./lobster workbuddy next
./lobster workbuddy open
```

macOS 使用提示：

- 如果提示“无法打开”或“已损坏”，通常是系统安全校验或隔离属性导致，可先在终端中直接运行二进制再根据系统提示处理
- 如果安装后仍然打不开应用，优先执行 `./lobster workbuddy doctor`

### Linux

适用产物示例：

- `lobster_<version>_linux_amd64.tar.gz`
- `lobster_<version>_linux_arm64.tar.gz`

建议步骤（Shell）：

```bash
# 1. 解压下载好的压缩包
tar -xzf lobster_v0.1.0_linux_amd64.tar.gz

# 2. 进入解压目录
cd lobster_v0.1.0_linux_amd64

# 3. 如有需要，补可执行权限
chmod +x lobster

# 4. 先查看帮助
./lobster help

# 5. 先做一次演练，不真正执行安装
./lobster workbuddy install --dry-run

# 6. 真正安装 WorkBuddy
./lobster workbuddy install

# 7. 安装完成后检查状态
./lobster workbuddy status

# 8. 查看下一步建议
./lobster workbuddy next
```

Linux 使用提示：

- 如果运行环境没有桌面，会影响 `open` 行为，但不影响安装和状态检测
- 如果安装完成后 shell 仍然找不到新命令，先重新打开终端，或执行 `./lobster workbuddy doctor` 查看 PATH 与安装痕迹

### 推荐使用顺序

无论你使用哪个平台，都建议按下面顺序操作：

```bash
lobster workbuddy install --dry-run
lobster workbuddy install
lobster workbuddy status
lobster workbuddy next
lobster workbuddy open
lobster workbuddy doctor
```

含义分别是：

- `install --dry-run`：先确认当前平台会执行什么安装策略
- `install`：真正执行安装
- `status`：检查是否已经安装成功
- `next`：获取安装后的下一步引导
- `open`：尝试打开 WorkBuddy
- `doctor`：遇到问题时输出诊断建议

## 本地开发

运行主命令：

```bash
go run ./cmd/lobster help
go run ./cmd/lobster tui
go run ./cmd/lobster list
go run ./cmd/lobster workbuddy help
go run ./cmd/lobster workbuddy tui
go run ./cmd/lobster workbuddy install --dry-run
```

编译：

```bash
go build ./...
```
