# Gobrew（中文说明）

[English](../README.md) | [中文](./README.zh-CN.md)

<p align="center">
  <img src="./logo.png" alt="Gobrew Logo" width="160" />
</p>

<p align="center">
  <strong>Homebrew 的桌面控制中心。</strong><br/>
  基于 Go + Wails + Vue，覆盖安装、更新、维护、服务管理的完整流程。
</p>

## 下载

- macOS 安装包与发布文件：[GitHub Releases](https://github.com/sunmking/Gobrew/releases)

## 为什么是 Gobrew

Homebrew CLI 很强大，但在日常批量操作时，命令行体验不一定高效。  
Gobrew 提供可视化桌面界面，底层仍然使用原生 Homebrew 能力。

- 更快搜索并安装 formula/cask
- 批量更新与清理，反馈更直观
- 在一个界面里管理 tap、services、Brewfile
- 保留原生 brew 工作流，提升交互体验

## 技术栈

- 后端：Go + Wails v3
- 前端：Vue 3 + TypeScript + Pinia + Vue Router + Vue I18n
- 样式：Tailwind + 统一自定义样式体系

## 快速开始

### 用户（下载安装）

1. 打开 [GitHub Releases](https://github.com/sunmking/Gobrew/releases)
2. 下载最新 `.dmg` 安装包
3. 在 macOS 上安装并启动 Gobrew

### 开发者（本地运行）

- macOS（主要目标平台）
- 已安装 Homebrew，且在 `PATH` 中可用
- Go `1.25+`
- Node.js + npm
- Wails v3 CLI
- Task（`go-task`，推荐）

```bash
# 安装前端依赖
npm --prefix frontend install

# 启动桌面开发模式
task dev
```

也可以直接运行 Wails：

```bash
wails3 dev -config ./build/config.yml -port 9245
```

## 核心使用场景

### Explore（搜索）
- 分类检索包结果
- 查看包详情与依赖关系
- 批量选择并执行安装

### Install（安装管理）
- 查看已安装 formula/cask
- 卸载、重装与信息查看
- 实时查看命令输出

### Maintain（维护）
- Tap 管理（`add/remove/list/info`）
- cleanup / autoremove 相关操作
- Brewfile 流程（`dump/list/check/cleanup/install/restore`）

### Services（服务）
- 查看后台服务状态
- 单个或批量 `start/stop/restart`
- 输出日志与完成事件追踪

## 构建与发布

```bash
task build
task package
```

## 常用命令

```bash
task dev            # 桌面开发模式
task run            # 运行已构建桌面应用
task build          # 构建桌面应用
task package        # 打包产物

task build:server   # 构建后端服务模式
task run:server     # 运行后端服务模式
task build:docker   # 构建 server 模式 docker 镜像
task run:docker     # 运行 docker 镜像
```

## 开发说明

### 项目结构

```text
Gobrew/
├── frontend/         # Vue 前端
├── services/         # Go 服务层（brew、tap、service、bundle）
├── build/            # Wails 配置与打包脚本
├── docs/             # 文档与视觉资源
├── main.go           # 应用启动与服务注册
└── Taskfile.yml      # 开发任务入口
```

### 运行模型

- 前端通过 Wails 绑定调用后端服务：
  - `BrewService`
  - `TapService`
  - `ServiceManager`
  - `BundleService`
- 长任务通过事件流回传 UI：
  - `brew-output`
  - `brew-complete`
- 内置 i18n（`en` / `zh`），可继续扩展

## 路线图

- 大结果集虚拟渲染进一步优化
- 更丰富的包元信息与依赖视图
- 自动化发布与签名构建
- 更完整的跨平台支持与文档

## 贡献

欢迎提交 Issue / PR。  
较大改动建议先开 Issue 对齐范围。

## 许可证

[MIT](../LICENSE)
