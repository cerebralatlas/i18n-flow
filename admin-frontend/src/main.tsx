import React from "react";
import ReactDOM from "react-dom/client";
import { ConfigProvider } from "antd";
import zhCN from "antd/lib/locale/zh_CN";
import App from "./App";
import "antd/dist/reset.css"; // 如果使用 Ant Design v5
import "@ant-design/v5-patch-for-react-19";
import "./index.css";

// 可能需要配置axios默认设置
import axios from "axios";
axios.defaults.baseURL = import.meta.env.VITE_API_URL || "";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <ConfigProvider locale={zhCN}>
      <App />
    </ConfigProvider>
  </React.StrictMode>
);
