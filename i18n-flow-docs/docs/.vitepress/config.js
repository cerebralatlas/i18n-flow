export default {
  title: "i18n Flow",
  description: "A modern i18n management system",

  themeConfig: {
    // Global theme configuration that applies to all locales
    logo: "/logo.png", // If you have a logo

    // Primary navigation and sidebar - this will be overridden by locale-specific settings
    nav: [
      { text: "Home", link: "/en/" },
      { text: "Guide", link: "/en/guide/getting-started" },
      { text: "API", link: "/en/api/backend-api" },
    ],

    sidebar: {
      "/en/guide/": [
        {
          text: "Introduction",
          items: [
            { text: "Getting Started", link: "/en/guide/getting-started" },
            { text: "Installation", link: "/en/guide/installation" },
            { text: "Frontend Guide", link: "/en/guide/frontend-guide" },
            { text: "Backend Guide", link: "/en/guide/backend-guide" },
            { text: "CLI Guide", link: "/en/guide/cli-guide" },
            { text: "Usage Tutorial", link: "/en/guide/usage" },
          ],
        },
      ],
      "/en/api/": [
        {
          text: "API Reference",
          items: [{ text: "Backend API", link: "/en/api/backend-api" }],
        },
      ],
      "/en/deployment/": [
        {
          text: "Deployment",
          items: [
            { text: "Docker Deployment", link: "/en/deployment/docker" },
            {
              text: "Kubernetes Deployment",
              link: "/en/deployment/kubernetes",
            },
          ],
        },
      ],
      "/en/faq/": [
        {
          text: "FAQ",
          items: [
            { text: "Troubleshooting", link: "/en/faq/troubleshooting" },
            { text: "Best Practices", link: "/en/faq/best-practices" },
          ],
        },
      ],

      // Chinese sidebar
      "/zh/guide/": [
        {
          text: "介绍",
          items: [
            { text: "开始使用", link: "/zh/guide/getting-started" },
            { text: "安装说明", link: "/zh/guide/installation" },
            { text: "前端使用指南", link: "/zh/guide/frontend-guide" },
            { text: "后端使用指南", link: "/zh/guide/backend-guide" },
            { text: "CLI工具指南", link: "/zh/guide/cli-guide" },
            { text: "使用教程", link: "/zh/guide/usage" },
          ],
        },
      ],
      "/zh/api/": [
        {
          text: "API参考",
          items: [{ text: "后端API", link: "/zh/api/backend-api" }],
        },
      ],
      "/zh/deployment/": [
        {
          text: "部署指南",
          items: [
            { text: "Docker部署", link: "/zh/deployment/docker" },
            { text: "Kubernetes部署", link: "/zh/deployment/kubernetes" },
            { text: "自定义配置", link: "/zh/deployment/custom-config" },
          ],
        },
      ],
      "/zh/faq/": [
        {
          text: "常见问题",
          items: [
            { text: "故障排除", link: "/zh/faq/troubleshooting" },
            { text: "最佳实践", link: "/zh/faq/best-practices" },
          ],
        },
      ],
    },
  },

  // Locales configuration
  locales: {
    root: {
      // Root redirects to English version
      // label: "English",
      lang: "en",
      link: "/en/",
    },
    en: {
      lang: "en",
      label: "English",
      link: "/en/",
      // English-specific theme config (overrides global)
      themeConfig: {
        nav: [
          { text: "Home", link: "/en/" },
          { text: "Guide", link: "/en/guide/getting-started" },
          { text: "API", link: "/en/api/backend-api" },
        ],
        // English language menu text
        langMenuLabel: "Language",
      },
    },
    zh: {
      lang: "zh",
      label: "简体中文",
      link: "/zh/",
      // Chinese-specific theme config (overrides global)
      themeConfig: {
        nav: [
          { text: "首页", link: "/zh/" },
          { text: "指南", link: "/zh/guide/getting-started" },
          { text: "API", link: "/zh/api/backend-api" },
        ],
        // Chinese language menu text
        langMenuLabel: "语言",
      },
    },
  },
};
