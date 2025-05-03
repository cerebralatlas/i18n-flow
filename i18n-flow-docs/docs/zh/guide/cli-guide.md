# CLI 工具指南

i18n Flow CLI 是一个功能强大的命令行工具，用于管理您的应用程序翻译，实现开发工作流与翻译管理系统的无缝集成。

## 特点

- 🔄 在项目和翻译服务器之间同步翻译
- 🔍 扫描源代码中的翻译键
- 📤 将新键推送到翻译服务器
- 📊 检查不同语言环境的翻译状态
- 📂 支持扁平和嵌套的翻译文件结构

## 安装

i18n Flow CLI 可以通过 npm 或 yarn 安装：

```bash
# 使用 npm 全局安装
npm install -g i18n-flow-cli

# 或使用 yarn 全局安装
yarn global add i18n-flow-cli

# 或使用 pnpm 全局安装
pnpm add -g i18n-flow-cli

# 或直接使用 npx
npx i18n-flow-cli <命令>
```

## 快速入门

### 1. 初始化项目

首先，您需要在项目中初始化 i18n Flow：

```bash
i18n-flow init
```

这将启动交互式设置向导，引导您完成配置过程，包括：

- 服务器 URL
- API 密钥
- 项目 ID
- 本地化目录
- 默认语言

完成后，配置将保存在项目根目录的 `.i18nflowrc.json` 文件中。

### 2. 从服务器同步翻译

初始化后，您可以从服务器同步翻译：

```bash
i18n-flow sync
```

这将从服务器下载所有翻译，并将它们保存到您配置的本地化目录中。

### 3. 推送新的翻译键

您可以扫描源代码并将新的翻译键推送到服务器：

```bash
i18n-flow push --scan
```

## 命令详解

### `init`

初始化 i18n Flow 配置。

```bash
i18n-flow init
```

这个交互式设置将引导您完成项目配置，包括：

- 服务器 URL
- API 密钥
- 项目 ID
- 本地化目录
- 默认语言

### `sync`

从服务器同步翻译到本地项目。

```bash
i18n-flow sync [选项]
```

选项：

- `-l, --locale <locales>` - 要同步的语言环境（逗号分隔）
- `-f, --force` - 强制覆盖本地翻译
- `-n, --nested` - 使用嵌套目录结构（每个语言一个文件夹，按键的第一部分分割）

示例：

```bash
# 同步所有语言
i18n-flow sync

# 同步特定语言
i18n-flow sync -l en,zh,ja

# 强制覆盖本地翻译
i18n-flow sync -f

# 使用嵌套结构同步
i18n-flow sync -n
```

### `push`

将翻译键推送到服务器。

```bash
i18n-flow push [选项]
```

选项：

- `-s, --scan` - 扫描源文件中的翻译键
- `-p, --path <patterns>` - 源文件的 Glob 模式（逗号分隔）
- `-d, --dry-run` - 显示键但不推送到服务器
- `-n, --nested` - 使用嵌套目录结构（从每个语言的文件夹读取）

示例：

```bash
# 推送已有的翻译键
i18n-flow push

# 扫描源代码并推送新键
i18n-flow push --scan

# 扫描特定文件
i18n-flow push -s -p "src/components/**/*.tsx,src/pages/**/*.tsx"

# 使用嵌套结构推送
i18n-flow push -n

# 仅显示扫描结果而不推送
i18n-flow push -s -d
```

### `status`

检查项目的翻译状态。

```bash
i18n-flow status
```

显示翻译状态信息，包括：

- 每种语言的键数量
- 服务器缺少的键
- 本地缺少的键
- 各语言的翻译完成率

## 文件结构

i18n Flow CLI 支持两种翻译文件结构：扁平结构和嵌套结构。

### 扁平结构（默认）

默认情况下，i18n Flow 将翻译组织为扁平结构：

```
locales/
  en.json
  zh.json
  fr.json
  ...
```

每个文件包含该语言的所有键：

```json
{
  "common.buttons.save": "保存",
  "common.buttons.cancel": "取消",
  "user.profile.title": "用户资料"
}
```

### 嵌套结构

使用 `-n, --nested` 选项时，i18n Flow 将翻译组织为嵌套结构：

```
locales/
  en/
    common.json
    user.json
  zh/
    common.json
    user.json
  ...
```

每个命名空间文件只包含相关键：

```json
// locales/zh/common.json
{
  "buttons": {
    "save": "保存",
    "cancel": "取消"
  }
}

// locales/zh/user.json
{
  "profile": {
    "title": "用户资料"
  }
}
```

## 配置

配置存储在项目根目录的 `.i18nflowrc.json` 文件中：

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

### 配置选项说明

- `serverUrl`: i18n Flow 服务器的 URL
- `apiKey`: API 密钥，用于 CLI 认证
- `projectId`: 项目 ID 或 slug
- `localesDir`: 翻译文件的本地目录
- `defaultLocale`: 默认语言代码
- `sourcePatterns`: 要扫描的源文件模式（支持 glob）
- `extractorPattern`: 用于提取翻译键的正则表达式

## 实践示例

### 示例 1: 在 React 项目中集成

假设您有一个 React 项目，使用 react-i18next 进行国际化：

1. 初始化 i18n Flow：

```bash
i18n-flow init
```

2. 配置提取模式以匹配 react-i18next 用法：

```json
{
  "extractorPattern": "(?:t|i18next\\.t|useTranslation\\(\\)[^}]*\\.)t)\\(['\"]([\\w\\.\\-]+)['\"]"
}
```

3. 从服务器同步最新翻译：

```bash
i18n-flow sync
```

4. 在开发过程中添加新的翻译键后：

```bash
i18n-flow push --scan
```

### 示例 2: 持续集成工作流

在 CI/CD 流程中，您可以自动化翻译同步：

```yaml
# 示例 GitHub Actions 工作流
name: Sync Translations

on:
  schedule:
    - cron: "0 0 * * *" # 每天运行一次

jobs:
  sync:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-node@v2
        with:
          node-version: "18"
      - run: npm install -g i18n-flow-cli
      - run: i18n-flow sync
      - name: Commit changes
        uses: stefanzweifel/git-auto-commit-action@v4
        with:
          commit_message: "chore: sync translations"
```

## 常见问题排查

### 认证失败

如果遇到认证错误：

1. 确认 API 密钥是否正确
2. 检查服务器 URL 是否正确
3. 确认服务器是否在运行

### 翻译键未被检测到

如果扫描没有检测到某些翻译键：

1. 确认提取模式（`extractorPattern`）是否匹配您的代码中使用的模式
2. 确认源文件模式（`sourcePatterns`）是否包含了所有相关文件

### 文件结构问题

如果翻译文件结构不符合预期：

1. 确认是否正确使用了 `-n` 选项
2. 检查 i18n 库的配置是否与 i18n Flow 的配置匹配
