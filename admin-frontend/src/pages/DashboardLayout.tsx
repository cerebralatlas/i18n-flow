import React, { useState } from "react";
import { Button, Layout, Menu } from "antd";
import {
  DashboardOutlined,
  TranslationOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  ProjectOutlined,
} from "@ant-design/icons";
import { useAuth } from "../contexts/AuthContext";
import { Outlet, useNavigate, useLocation } from "react-router-dom";

const { Header, Sider, Content } = Layout;

const DashboardLayout: React.FC = () => {
  const { user, logout } = useAuth();
  const [collapsed, setCollapsed] = useState(false);
  const navigate = useNavigate();
  const location = useLocation();

  // Determine the selected menu item based on the current path
  const getSelectedKey = () => {
    const path = location.pathname;
    if (path.includes("/projects")) return "2";
    if (path.includes("/translations")) return "3";
    return "1"; // default dashboard
  };

  return (
    <Layout style={{ minHeight: "100vh" }}>
      <Sider
        trigger={null}
        collapsible
        collapsed={collapsed}
        theme="light"
        className="shadow-md"
      >
        <div className="p-4 flex justify-center items-center h-16 border-b">
          <h1
            className={`text-xl font-bold transition-opacity duration-300 ${
              collapsed ? "opacity-0 w-0" : "opacity-100"
            }`}
          >
            i18n-flow
          </h1>
        </div>
        <Menu
          mode="inline"
          selectedKeys={[getSelectedKey()]}
          onClick={({ key }) => {
            switch (key) {
              case "1":
                navigate("/dashboard");
                break;
              case "2":
                navigate("/projects");
                break;
              case "3":
                navigate("/translations");
                break;
            }
          }}
          items={[
            {
              key: "1",
              icon: <DashboardOutlined />,
              label: "Dashboard",
            },
            {
              key: "2",
              icon: <ProjectOutlined />,
              label: "Project Management",
            },
            {
              key: "3",
              icon: <TranslationOutlined />,
              label: "Translation Management",
            },
          ]}
        />
      </Sider>
      <Layout>
        <Header className="bg-white p-0 flex justify-between items-center px-4 shadow-sm">
          <Button
            type="text"
            icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
            onClick={() => setCollapsed(!collapsed)}
          />
          <div className="flex items-center gap-4">
            <span>Welcome, {user?.username}</span>
            <Button onClick={logout}>Logout</Button>
          </div>
        </Header>
        <Content className="m-6">
          <Outlet />
        </Content>
      </Layout>
    </Layout>
  );
};

export default DashboardLayout;
