# i18n-flow CLI

一个强大的 CLI 工具，轻松管理您应用程序的翻译内容。

[![npm version](https://img.shields.io/npm/v/i18n-flow-cli.svg)](https://www.npmjs.com/package/i18n-flow-cli)

简体中文 | [English](README.md)

## 功能特点

- 🔄 在项目和翻译服务器之间同步翻译
- 🔍 扫描源代码中的翻译键
- 📤 将新键推送到翻译服务器
- 📊 检查不同语言环境的翻译状态
- 📂 支持扁平和嵌套的翻译文件结构

## 安装

```bash
# 全局安装
npm install -g i18n-flow-cli

# 或使用npx
npx i18n-flow-cli <命令>
```

## 快速开始

1. 在项目中初始化 i18n-flow：

```bash
i18n-flow init
```

2. 从服务器同步翻译：

```bash
i18n-flow sync
```

3. 扫描翻译键并推送到服务器：

```bash
i18n-flow push --scan
```

## 命令

### `init`

在项目中初始化 i18n-flow。

```bash
i18n-flow init
```

这个交互式设置将引导您完成使用 i18n-flow 配置项目的过程，包括：

- 服务器 URL
- API 密钥
- 项目 ID
- 语言文件目录
- 默认语言环境

### `sync`

将翻译从服务器同步到本地项目。

```bash
i18n-flow sync [选项]
```

选项：

- `-l, --locale <locales>` - 要同步的语言环境的逗号分隔列表
- `-f, --force` - 强制覆盖本地翻译
- `-n, --nested` - 使用嵌套目录结构（每个语言环境一个文件夹，按键的第一部分拆分）

### `push`

将翻译键推送到服务器。

```bash
i18n-flow push [选项]
```

选项：

- `-s, --scan` - 扫描源文件中的翻译键
- `-p, --path <patterns>` - 源文件的 Glob 模式（逗号分隔）
- `-d, --dry-run` - 显示键而不推送到服务器
- `-n, --nested` - 使用嵌套目录结构（从每个语言环境的文件夹读取）

### `status`

检查项目的翻译状态。

```bash
i18n-flow status
```

显示翻译的状态，包括：

- 每个语言环境的键数量
- 服务器缺少的键
- 服务器上没有的额外键

## 文件结构

### 传统（扁平）结构

默认情况下，i18n-flow 以扁平结构组织翻译：

```
locales/
  en.json
  fr.json
  zh.json
  ...
```

每个文件包含该语言环境的所有键：

```json
{
  "account.info": "账户信息",
  "account.settings": "账户设置",
  "dashboard.title": "仪表板"
}
```

### 嵌套结构

使用`-n, --nested`选项，i18n-flow 以嵌套结构组织翻译：

```
locales/
  en/
    account.json
    dashboard.json
  fr/
    account.json
    dashboard.json
  ...
```

每个命名空间文件只包含相关键：

```json
// locales/en/account.json
{
  "info": "账户信息",
  "settings": "账户设置"
}

// locales/en/dashboard.json
{
  "title": "仪表板"
}
```

## 配置

配置存储在项目根目录的`.i18nflowrc.json`中：

```json
{
  "serverUrl": "https://your-i18n-server.com",
  "apiKey": "your-api-key",
  "projectId": "your-project-id",
  "localesDir": "./src/locales",
  "defaultLocale": "en",
  "sourcePatterns": [
    "src/**/*.{js,jsx,ts,tsx}",
    "!src/**/*.{spec,test}.{js,jsx,ts,tsx}",
    "!**/node_modules/**"
  ],
  "extractorPattern": "(?:t|i18n\\.t)\\(['\"]([\\w\\.\\-]+)['\"]"
}
```

## 使用示例

### 初始化新项目

```bash
i18n-flow init
```

### 使用嵌套结构同步翻译

```bash
i18n-flow sync --nested
```

### 从嵌套结构推送键

```bash
i18n-flow push --nested
```

### 扫描特定文件以查找键

```bash
i18n-flow push --scan --path "src/components/**/*.tsx,src/pages/**/*.tsx"
```

### 检查翻译状态

```bash
i18n-flow status
```
