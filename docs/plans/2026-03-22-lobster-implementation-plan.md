# Lobster 实施计划

日期：2026-03-22

## 1. 第一阶段目标

交付一个可运行的跨平台 CLI 原型，用于安装腾讯 WorkBuddy，并提供基础状态查询、启动与诊断能力。

## 2. 分阶段实施

### 阶段一：项目初始化

目标：

- 初始化 Go CLI 项目
- 明确包结构
- 设定命令入口

建议输出：

- `go.mod`
- CLI 主入口
- `install/status/open/doctor/list` 基础命令骨架

### 阶段二：平台检测与安装编排

目标：

- 识别 Windows、macOS、Linux
- 为 WorkBuddy 选择官方安装路径
- 执行安装并收集结果

建议输出：

- 平台检测模块
- 官方安装执行器
- 安装结果对象

### 阶段三：安装结果检测

目标：

- 判断 WorkBuddy 是否安装成功
- 检查 PATH 与应用路径
- 形成统一的状态输出

建议输出：

- 检测器模块
- 状态对象
- 用户可读的安装结论

### 阶段四：启动与下一步建议

目标：

- 拉起 WorkBuddy
- 输出下一步引导
- 为后续绑定流程留扩展点

建议输出：

- 启动器模块
- 建议器模块
- `next` 或安装成功提示逻辑

### 阶段五：诊断与异常处理

目标：

- 对常见失败给出稳定提示
- 建立统一错误分类

建议输出：

- 诊断规则
- 错误码或错误类型
- 常见问题映射表

## 3. 推荐包结构

```text
cmd/lobster
internal/platform
internal/installer
internal/products/workbuddy
internal/detector
internal/launcher
internal/advisor
internal/cli
```

## 4. 关键设计约束

- 不复制官方安装脚本逻辑，优先调用官方入口
- 不解析脆弱的安装输出文本作为唯一成功依据
- 不依赖 WorkBuddy 内部私有文件结构
- 安装逻辑必须幂等
- 默认输出要适合新手阅读

## 5. 第一版建议命令

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
```

## 6. 测试策略

### 单元测试

- 平台识别逻辑
- 命令参数解析
- 错误映射逻辑
- 状态建议逻辑

### 集成测试

- 模拟不同平台命令选择
- 模拟安装成功/失败场景
- 模拟 PATH 未生效与权限问题

### 手工验证

- macOS 真机验证
- Windows 真机验证
- Linux 桌面环境验证

## 7. 风险点

- 官方安装器行为变化
- 平台差异导致启动路径不一致
- Windows 执行策略差异
- Linux 发行版桌面能力差异

## 8. 当前优先级建议

推荐优先级：

1. `install workbuddy`
2. `status workbuddy`
3. `open workbuddy`
4. `doctor workbuddy`
5. `list`

## 9. 实施原则

- 先把 WorkBuddy 跑通
- 再抽象多产品支持
- 不为未来扩展过度设计
- 但命令结构必须为未来预留空间

## 10. 当前进展（2026-03-23）

目前仓库已经完成以下落地项：

- 已建立 `installer.Result` 安装闭环结果模型，包含安装前状态、安装后状态与安装结果结论
- `install workbuddy` 已接入安装前检测、重复安装跳过、安装后复检与校验失败提示
- `status workbuddy` 已能区分“命令可用”“仅有安装痕迹”“未检测到安装”
- `doctor workbuddy` 已能输出命令状态、路径证据、桌面环境提示与下一步建议
- 已补充 `advisor`、`cli`、`detector`、`platform`、`products` 的基础单元测试

下一阶段建议优先处理：

1. 继续细化 `doctor` 的错误分类与规则覆盖
2. 为 `cli` 和 `installer` 增加更稳定的依赖注入与测试入口
3. 逐步为多产品扩展保留更清晰的命令建议生成方式
