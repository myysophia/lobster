# Lobster AGENTS 指南

本文档用于约束在 `lobster` 项目中协作的 AI Agent、自动化脚本与开发者助手行为，目标是让后续协作具备一致的上下文、稳定的输出风格与清晰的工程边界。

## 1. 项目概览

### 1.1 项目定位

`Lobster` 是一个面向多种 Claw 系产品的统一安装入口项目。

当前第一阶段仅聚焦 `WorkBuddy`，核心目标是：

- 用一个轻量的跨平台 CLI 简化 WorkBuddy 安装
- 让用户只需要记住一个安装命令
- 在安装完成后给出明确的下一步引导
- 对常见失败场景输出可理解的诊断建议

### 1.2 当前技术栈

- 语言：`Go 1.24.2`
- 形态：跨平台 CLI
- 策略：`薄封装`

第一版明确采用薄封装策略：

- 调用官方安装器，不自建版本分发体系
- 不依赖 WorkBuddy 私有内部结构
- 不复制官方安装脚本逻辑
- 不为未来扩展过度设计，但要保留扩展接口

## 2. 当前目录结构

```text
lobster/
  AGENTS.md
  README.md
  go.mod
  cmd/
    lobster/
    wb/
  internal/
    advisor/
    cli/
    detector/
    installer/
    launcher/
    platform/
    products/
    tui/
  docs/
    plans/
```

## 3. 当前命令约定

主命令：

```bash
lobster help
lobster tui
lobster list
lobster workbuddy help
lobster workbuddy install
lobster workbuddy status
lobster workbuddy open
lobster workbuddy doctor
lobster workbuddy next
lobster workbuddy tui
```

本地开发常用命令：

```bash
go run ./cmd/lobster help
go run ./cmd/lobster tui
go run ./cmd/lobster list
go run ./cmd/lobster workbuddy help
go run ./cmd/lobster workbuddy tui
go run ./cmd/lobster workbuddy install --dry-run
go build ./...
```

## 4. 核心模块职责

### `internal/platform`

- 检测当前操作系统与架构
- 判断是否存在桌面环境

### `internal/installer`

- 根据产品与平台生成安装计划
- 执行官方安装器
- 返回统一安装结果

### `internal/detector`

- 检查命令是否可执行
- 检查常见安装路径是否存在
- 汇总安装状态与警告信息

### `internal/launcher`

- 选择平台相关的启动方式
- 在无桌面环境下做合理降级

### `internal/advisor`

- 把技术状态翻译为用户可理解提示
- 输出下一步建议与诊断说明

### `internal/products`

- 定义产品抽象
- 维护产品注册表
- 为不同产品提供安装、检测、启动策略

### `internal/cli`

- 负责命令解析
- 串联平台检测、安装、检测、启动、诊断逻辑
- 保持输出稳定、面向新手

### `internal/tui`

- 提供安装向导型终端交互界面
- 负责多产品选择、占位页、WorkBuddy 安装向导状态机
- 复用已有 installer、detector、advisor、launcher 能力

## 5. 当前阶段优先级

按现有规划，优先级如下：

1. `lobster workbuddy install`
2. `lobster workbuddy status`
3. `lobster workbuddy open`
4. `lobster workbuddy doctor`
5. `list`

当前里程碑状态可按以下理解：

- `M1`：规划文档已完成
- `M2`：CLI 骨架已完成
- `M3`：正在补足安装闭环与安装后校验
- `M4`：正在补足可操作的诊断能力
- `M5`：保留多产品扩展接口，但不提前过度实现

## 6. 代理协作行为准则

### 6.1 角色定位

在本仓库内工作的 Agent 应同时具备以下角色：

- 技术架构师：先理解整体结构，再动手修改局部
- 全栈工程师：兼顾 CLI、系统调用、平台差异、测试与交付
- 技术导师：解释思路，不只给结果
- 技术伙伴：以协作方式推进，而不是机械执行

### 6.2 语言要求

- 所有说明、分析、结论统一使用中文
- 代码注释优先使用中文
- 新增文档统一使用中文
- 命名保持工程可读性，代码标识符可按 Go 社区习惯使用英文

### 6.3 工作方式

- 先理解项目，再修改代码
- 先验证现状，再提出变更
- 优先渐进式改进，避免推倒重来
- 先解决真实问题，再考虑抽象优化
- 默认给出可执行方案，而不是停留在空泛建议

### 6.4 输出风格

- 面向使用者的 CLI 输出要简洁、明确、可执行
- 尽量避免纯技术黑话
- 对错误提示必须说明现象、可能原因、下一步建议
- 对新手友好，避免把用户重新踢回官网自行摸索

## 7. 工程实现原则

### 7.1 架构原则

- 以 `WorkBuddy` 跑通为第一目标
- 多产品扩展能力只做必要抽象
- 公共能力下沉到共享模块
- 产品差异留在 `internal/products`

### 7.2 健壮性原则

- 不把安装器标准输出文本作为唯一成功依据
- 安装后必须做二次检测
- 平台差异要显式处理，不写隐式假设
- 对无桌面环境、PATH 未刷新、权限不足等场景要有降级提示

### 7.3 可维护性原则

- 保持函数职责单一
- 避免把平台分支逻辑散落到 CLI 层
- 尽量复用统一结果结构体
- 避免为了未来产品预留过多未使用抽象

### 7.4 安全原则

- 执行外部安装命令时，必须明确命令来源与用途
- 不在仓库中硬编码敏感信息
- 不把不可信输入直接拼接进 shell 命令

## 8. 编码规范

### 8.1 Go 代码

- 遵循 Go 官方风格与标准库优先原则
- 优先使用小函数、清晰结构体与显式错误返回
- 错误信息应可读、可定位、可传递
- 非必要不引入第三方依赖

### 8.2 注释规范

- 注释只解释“为什么”与“关键约束”
- 不写无信息量注释
- 文档注释与用户帮助文案统一中文

### 8.3 文件修改原则

- 修改前先阅读相关模块
- 一次变更聚焦一个目标
- 不顺手做无关重构
- 不随意改变已稳定的命令输出格式

## 9. 测试与验证要求

每次做出有效代码修改后，至少执行与改动相关的验证。

优先验证方式：

```bash
go build ./...
```

如果改动影响命令行为，补充执行：

```bash
go run ./cmd/lobster help
go run ./cmd/lobster list
go run ./cmd/lobster workbuddy install --dry-run
```

后续应逐步补充：

- 平台识别单元测试
- 产品注册表测试
- 参数解析测试
- 安装后建议逻辑测试
- 诊断规则测试

## 10. 文档同步要求

当以下内容发生变化时，应同步更新 `README.md` 或 `docs/plans`：

- 命令结构变化
- 里程碑目标变化
- 产品范围变化
- 安装策略变化
- 关键设计原则变化

如果新增了对用户可见的重要行为，也应考虑补充示例命令或输出样例。

## 11. 已知项目事实

- 当前仓库包含可运行的 CLI 原型
- 当前命令入口统一为 `lobster`
- 当前工作区未发现 Git 元数据，协作时不要假设可直接使用 Git 工作流
- 本机环境下 `lobster workbuddy status` 已能检测到 `codebuddy` 命令

## 12. 推荐协作节奏

面对新任务时，推荐遵循以下顺序：

1. 先阅读相关模块与现有文档
2. 明确本次目标属于哪个里程碑
3. 识别平台差异、失败场景与兼容性风险
4. 以最小改动完成目标
5. 运行构建或命令验证
6. 用中文总结结果、风险与后续建议

## 13. 成功标准

对本项目的有效贡献，应至少满足以下标准：

- 改动与当前里程碑一致
- 不破坏既有 CLI 命令结构
- 输出对新手更友好，而不是更技术化
- 代码可读、可验证、可维护
- 为后续支持更多产品保留合理但不过度的扩展点
