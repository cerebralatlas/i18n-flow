# 使用教程

本教程将引导您完成在真实项目中使用 i18n Flow 的完整流程，包括项目设置、翻译管理和开发集成。

## 准备工作

在开始使用 i18n Flow 之前，请确保完成以下准备工作：

1. 已成功部署 i18n Flow 系统（后端和前端）
2. 已安装 i18n Flow CLI 工具
3. 已准备要进行国际化的应用程序项目

## 基本工作流程

i18n Flow 的典型工作流程如下：

1. 在管理面板创建项目
2. 定义支持的语言
3. 使用 CLI 扫描源代码中的翻译键
4. 将翻译键推送到服务器
5. 在管理面板中添加翻译
6. 将翻译同步回开发项目
7. 在应用程序中集成翻译内容

## 步骤 1: 创建项目

首先，您需要在 i18n Flow 管理界面创建项目：

1. 登录管理界面
2. 点击"项目"菜单
3. 点击"创建项目"按钮
4. 填写项目信息：
   - **项目名称**：您的应用名称，例如 "我的在线商店"
   - **项目描述**：简要描述，例如 "电子商务平台的国际化内容"
   - **项目标识符(Slug)**：用于 API 和 CLI 调用的唯一标识符，例如 "my-online-store"
5. 点击"保存"按钮

## 步骤 2: 添加语言

接下来，添加您的应用程序需要支持的语言：

1. 在侧边栏菜单中点击"语言"
2. 点击"添加语言"按钮
3. 添加您需要支持的每种语言，例如：
   - 英语 (en)
   - 简体中文 (zh-CN)
   - 日语 (ja)
   - 西班牙语 (es)
4. 对于每种语言，填写以下信息：
   - **语言名称**：例如 "简体中文"
   - **语言代码**：例如 "zh-CN"
   - **本地名称**：例如 "中文"
   - **RTL 支持**：对于从右到左书写的语言（如阿拉伯语或希伯来语）选择是，否则选择否

## 步骤 3: 初始化项目中的 CLI

在您的开发项目中，初始化 i18n Flow CLI：

1. 打开终端并导航到项目目录
2. 运行初始化命令：

```bash
i18n-flow init
```

3. 按照提示配置 CLI：
   - **服务器 URL**：i18n Flow 服务器的地址，例如 "http://your-i18n-flow-server.com"
   - **API 密钥**：在 i18n Flow 后端配置中设置的 CLI API 密钥
   - **项目 ID**：您刚刚创建的项目 slug，例如 "my-online-store"
   - **本地化目录**：翻译文件的存放目录，例如 "./src/locales"
   - **默认语言**：主要开发语言，例如 "en"

## 步骤 4: 扫描和推送翻译键

现在扫描您的源代码，找出所有需要翻译的文本，并将它们推送到服务器：

```bash
i18n-flow push --scan
```

这将：

1. 扫描您的源代码
2. 提取所有匹配翻译模式的键
3. 将这些键推送到 i18n Flow 服务器

::: tip 提示
CLI 使用正则表达式 `extractorPattern` 来匹配翻译函数。默认模式是 `(?:t|i18n\.t)\(['"]([\\w\\.\\-]+)['"]`，它能匹配如 `t('common.hello')` 或 `i18n.t('user.greeting')` 这样的调用。如果您的项目使用不同的翻译函数，请相应地调整此模式。
:::

## 步骤 5: 添加翻译内容

键推送到服务器后，您可以添加翻译内容：

1. 在管理界面中，导航到您的项目
2. 切换到"翻译"选项卡
3. 您将看到所有从代码中提取的键
4. 为每个键添加不同语言的翻译内容
5. 您可以逐个编辑，也可以使用批量导入功能

### 批量导入翻译

对于大型项目，您可以使用批量导入功能：

1. 点击"导入"按钮
2. 选择文件格式（JSON、Excel、CSV）
3. 上传您的翻译文件或粘贴内容
4. 选择要更新的语言
5. 点击"导入"开始处理

### 导出翻译

您也可以导出翻译以便共享或备份：

1. 点击"导出"按钮
2. 选择要导出的语言
3. 选择文件格式
4. 点击"导出"下载文件

## 步骤 6: 同步翻译到项目

翻译内容准备好后，将其同步到您的开发项目：

```bash
i18n-flow sync
```

这将从服务器下载所有翻译，并将它们保存到您配置的本地化目录中，格式为 JSON 文件。

::: tip 文件结构
i18n Flow CLI 支持两种文件结构：

- **扁平结构**（默认）：每种语言一个文件，包含所有键
- **嵌套结构**：使用 `-n` 选项，按命名空间组织键
  :::

## 步骤 7: 集成到您的应用程序

最后，将翻译集成到您的应用程序中。这取决于您使用的前端框架和 i18n 库。

### React 项目示例（使用 react-i18next）

1. 安装 i18next 和 react-i18next：

```bash
npm install i18next react-i18next i18next-http-backend i18next-browser-languagedetector
```

2. 创建 i18n 配置文件（src/i18n.js）：

```javascript
import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import Backend from "i18next-http-backend";
import LanguageDetector from "i18next-browser-languagedetector";

// 导入翻译文件
import enTranslation from "./locales/en.json";
import zhTranslation from "./locales/zh-CN.json";

i18n
  .use(Backend)
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    resources: {
      en: { translation: enTranslation },
      "zh-CN": { translation: zhTranslation },
    },
    fallbackLng: "en",
    debug: process.env.NODE_ENV === "development",
    interpolation: {
      escapeValue: false,
    },
  });

export default i18n;
```

3. 在 App.js 中初始化 i18n：

```jsx
import React from "react";
import { useTranslation } from "react-i18next";
import "./i18n";

function App() {
  const { t } = useTranslation();

  return (
    <div className="App">
      <h1>{t("app.title")}</h1>
      <p>{t("app.greeting")}</p>
      <button>{t("common.buttons.save")}</button>
    </div>
  );
}

export default App;
```

## 持续集成流程

在实际开发中，您可以将 i18n Flow 集成到您的 CI/CD 流程中：

### 开发阶段

1. 开发人员在代码中添加新的翻译键
2. 开发人员运行 `i18n-flow push --scan` 将新键推送到服务器
3. 翻译人员通过管理界面添加翻译内容

### 构建阶段

1. CI 系统在构建前运行 `i18n-flow sync` 获取最新翻译
2. 应用程序使用最新翻译进行构建

### 示例 GitHub Actions 工作流

```yaml
name: Build with translations

on:
  push:
    branches: [main]

jobs:
  build:
    runs-on: ubuntu-latest

    steps:
      - uses: actions/checkout@v2

      - name: Setup Node.js
        uses: actions/setup-node@v2
        with:
          node-version: "18"

      - name: Install dependencies
        run: npm ci

      - name: Install i18n Flow CLI
        run: npm install -g i18n-flow-cli

      - name: Sync translations
        run: i18n-flow sync
        env:
          I18N_FLOW_API_KEY: ${{ secrets.I18N_FLOW_API_KEY }}

      - name: Build application
        run: npm run build
```

## 高级用法

### 翻译基于命名空间的组织

对于大型项目，您可能希望将翻译组织为命名空间：

1. 初始化时配置嵌套结构：

```bash
i18n-flow init
# 在配置向导中选择嵌套结构
```

2. 或者，在已配置的项目中使用嵌套选项：

```bash
i18n-flow sync -n
i18n-flow push -n
```

### 在代码重构时处理键变更

当您重构代码并更改翻译键时：

1. 在管理界面中，使用"批量编辑"功能重命名键
2. 更新代码中的键引用
3. 运行 `i18n-flow push --scan` 更新服务器上的键
4. 运行 `i18n-flow sync` 同步更改回项目

## 故障排除

### 键未被扫描到

如果某些键未被检测到：

1. 检查键格式是否与提取模式匹配
2. 确认文件路径是否包含在扫描模式中
3. 调整 `.i18nflowrc.json` 中的 `extractorPattern` 以匹配您的翻译函数调用方式

### 同步错误

如果遇到同步错误：

1. 确认 API 密钥和服务器 URL 是否正确
2. 检查网络连接
3. 验证项目 ID 是否存在于服务器上

## 最佳实践

1. **使用描述性键**：使用有意义的键名，如 `user.profile.title` 而不是 `title1`
2. **添加上下文**：在管理界面为每个键添加描述，帮助翻译人员理解上下文
3. **定期同步**：定期同步翻译，确保开发环境中的翻译是最新的
4. **自动化流程**：将翻译同步集成到您的 CI/CD 流程中
5. **使用默认文本**：在代码中提供默认文本，以便在翻译缺失时显示合理的内容
