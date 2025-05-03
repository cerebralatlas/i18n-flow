# 最佳实践

本指南提供了使用 i18n Flow 的推荐最佳实践，以最大限度地提高效率和可维护性。

## 项目组织

### 翻译键结构

- **使用层次命名空间**：使用命名空间组织键（例如，`common.buttons.save`，`user.profile.title`）
- **保持一致性**：在整个项目中遵循相同的命名模式
- **使用有意义的键**：选择描述性名称而不是抽象标识符

良好键结构的示例：

```json
{
  "common": {
    "buttons": {
      "save": "保存",
      "cancel": "取消",
      "delete": "删除"
    },
    "messages": {
      "success": "操作成功",
      "error": "发生错误"
    }
  },
  "user": {
    "profile": {
      "title": "用户资料",
      "email": "电子邮箱"
    }
  }
}
```

### 项目结构

- 考虑将大型应用程序拆分为多个 i18n Flow 项目：
  - 每个大型功能区域一个项目
  - 不同产品线使用单独的项目
  - 共享组件使用共享项目

## 翻译工作流

### 高效的键管理

1. **提前规划键结构**：在实施前设计翻译键层次结构
2. **避免更改键**：一旦使用，更改键会扰乱翻译
3. **添加上下文注释**：为翻译人员提供清晰的上下文
4. **一致使用占位符**：标准化占位符格式（例如，`{{name}}` 或 `{name}`）

### 开发流程

1. **定期扫描和推送**：将键扫描集成到开发工作流程中
2. **推送前审查**：在推送到服务器前验证新键
3. **使用持续集成**：设置 CI 任务自动同步翻译
4. **跟踪完成率**：监控每种语言的翻译完成情况

## 技术集成

### 代码集成

- **使用一致的翻译函数**：标准化在代码中访问翻译的方式
- **实现后备策略**：为缺失的翻译制定计划
- **添加类型检查**：使用 TypeScript 或 PropTypes 为翻译键提供类型安全
- **延迟加载翻译**：仅在需要时加载翻译，提高性能

良好 React 实现的示例：

```jsx
// 使用 react-i18next
function Button({ type = "save", onClick }) {
  const { t } = useTranslation();

  return (
    <button onClick={onClick} className={`btn btn-${type}`}>
      {t(`common.buttons.${type}`)}
    </button>
  );
}
```

### CLI 配置

- **微调提取器模式**：自定义正则表达式模式以匹配代码风格
- **排除测试文件**：自定义源模式以排除测试和模拟文件
- **大型项目使用嵌套结构**：当键数量增长时切换到嵌套结构

优化的 CLI 配置示例：

```json
{
  "serverUrl": "https://i18n-flow.example.com",
  "apiKey": "your-api-key",
  "projectId": "my-project",
  "localesDir": "./src/i18n/locales",
  "defaultLocale": "en",
  "sourcePatterns": [
    "src/**/*.{js,jsx,ts,tsx}",
    "!src/**/*.{spec,test}.{js,jsx,ts,tsx}",
    "!src/mocks/**",
    "!**/node_modules/**"
  ],
  "extractorPattern": "(?:t|i18n\\.t|useTranslation\\(\\)[^}]*\\.)t)\\(['\"]([\\w\\.\\-]+)['\"]"
}
```

## 翻译内容

### 为翻译人员写作

- **保持文本简单**：使用清晰、简洁的语言
- **避免习语和俚语**：它们难以准确翻译
- **提供上下文**：为潜在模糊的文本添加描述
- **包含截图**：视觉帮助翻译人员理解上下文

### 质量控制

- **审核翻译**：让母语人士审核翻译
- **在上下文中测试**：在实际 UI 中检查翻译
- **注意字符串长度**：意识到翻译可能比英文更长
- **验证变量占位符**：确保所有占位符在翻译中都得到保留

## 性能优化

### 前端性能

- **按路由拆分翻译**：仅加载当前路由的翻译
- **缓存翻译**：实现客户端缓存
- **优化包大小**：从生产构建中删除未使用的翻译

### 后端性能

- **实现缓存**：缓存 API 响应以减少数据库负载
- **使用批量操作**：优先使用批量 API 调用而不是单独请求
- **监控 API 使用情况**：跟踪 API 性能以识别瓶颈

## CI/CD 集成

### 自动化工作流

- **自动化翻译同步**：在 CI 管道中包含翻译同步
- **验证翻译完整性**：如果关键翻译缺失则使构建失败
- **生成翻译报告**：在构建过程中创建完成报告

GitHub Actions 工作流示例：

```yaml
name: i18n 工作流

on:
  push:
    branches: [main]
  pull_request:
    branches: [main]

jobs:
  i18n-check:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-node@v2
        with:
          node-version: "18"

      - name: 安装依赖
        run: npm ci

      - name: 安装 i18n Flow CLI
        run: npm install -g i18n-flow-cli

      - name: 推送新键
        run: i18n-flow push --scan

      - name: 检查翻译状态
        run: i18n-flow status --min-completion=80

      - name: 同步翻译
        run: i18n-flow sync
```

## 版本策略

### 管理版本

- **为翻译添加版本**：为主要应用程序发布使用版本标签
- **归档旧键**：不要删除旧键，将它们标记为已弃用
- **计划升级**：制定主要版本之间迁移翻译的策略

通过遵循这些最佳实践，您可以使用 i18n Flow 维护高效、可扩展和高质量的翻译工作流。
