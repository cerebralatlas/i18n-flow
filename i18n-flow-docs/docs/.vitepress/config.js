export default {
  title: "i18n Flow",
  description: "A modern i18n management system",
  locales: {
    root: {
      label: "English",
      lang: "en",
      link: "/en/",
    },
    en: {
      label: "English",
      lang: "en",
      link: "/en/",
      themeConfig: {
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
        },
      },
    },
  },
};
