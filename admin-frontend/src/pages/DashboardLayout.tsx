import React, { useState } from "react";
import { Button, Layout, Menu, Avatar, Dropdown } from "antd";
import {
  DashboardOutlined,
  TranslationOutlined,
  MenuFoldOutlined,
  MenuUnfoldOutlined,
  ProjectOutlined,
  UserOutlined,
  LogoutOutlined,
  GlobalOutlined,
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

  const userMenuItems = [
    {
      key: "profile",
      icon: <UserOutlined />,
      label: "Profile",
    },
    {
      type: "divider",
    },
    {
      key: "logout",
      icon: <LogoutOutlined />,
      label: "Logout",
      onClick: logout,
    },
  ];

  return (
    <Layout style={{ minHeight: "100vh" }}>
      <Sider
        trigger={null}
        collapsible
        collapsed={collapsed}
        width={220}
        style={{
          background: "#0f172a",
          boxShadow: "rgba(0, 0, 0, 0.1) 0px 4px 12px",
          borderRight: "none",
          zIndex: 10,
        }}
      >
        <div className="px-4 py-5 flex items-center h-16 border-b border-gray-800">
          <div className="flex items-center gap-2">
            <div className="flex justify-center items-center w-8 h-8 rounded bg-white">
              <GlobalOutlined className="text-white" />
            </div>
            {!collapsed && (
              <h1 className="text-lg font-medium text-white m-0 ml-2">
                i18n-flow
              </h1>
            )}
          </div>
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
          style={{
            background: "transparent",
            borderRight: "none",
            marginTop: "12px",
          }}
          theme="dark"
        />

        {!collapsed && (
          <div className="absolute bottom-4 left-0 right-0 px-4">
            <div className="rounded border border-gray-800 p-3">
              <div className="flex items-center">
                <Avatar
                  size={32}
                  icon={<UserOutlined />}
                  style={{ backgroundColor: "#1677ff" }}
                />
                <div className="ml-3">
                  <span className="text-white text-sm block">
                    {user?.username}
                  </span>
                  <span className="text-gray-400 text-xs">Admin</span>
                </div>
              </div>
            </div>
          </div>
        )}
      </Sider>

      <Layout>
        <Header
          className="px-6 flex justify-between items-center"
          style={{
            background: "#ffffff",
            height: "64px",
            padding: "0 16px",
            boxShadow: "0 1px 2px rgba(0, 0, 0, 0.03)",
            position: "sticky",
            top: 0,
            zIndex: 9,
            borderBottom: "1px solid #f0f0f2",
          }}
        >
          <div className="flex items-center">
            <Button
              type="text"
              icon={collapsed ? <MenuUnfoldOutlined /> : <MenuFoldOutlined />}
              onClick={() => setCollapsed(!collapsed)}
              className="mr-3"
              style={{ fontSize: "16px" }}
            />
            <span className="font-medium text-gray-800">
              {location.pathname.includes("/projects") && "Project Management"}
              {location.pathname.includes("/translations") &&
                "Translation Management"}
              {location.pathname === "/dashboard" && "Dashboard"}
            </span>
          </div>

          <div className="flex items-center space-x-3">
            <Dropdown
              menu={{ items: userMenuItems as any }}
              placement="bottomRight"
              arrow
            >
              <div className="flex items-center cursor-pointer px-1 py-1 rounded hover:bg-gray-50">
                <Avatar
                  size={28}
                  icon={<UserOutlined />}
                  style={{ backgroundColor: "#1677ff" }}
                />
                <span className="ml-2 text-sm">{user?.username}</span>
              </div>
            </Dropdown>
          </div>
        </Header>

        <Content
          style={{
            padding: "24px",
            background: "#f5f5f5",
          }}
        >
          <div
            style={{
              padding: "24px",
              backgroundColor: "white",
              boxShadow: "0 1px 2px rgba(0, 0, 0, 0.03)",
              borderRadius: "2px",
            }}
          >
            <Outlet />
          </div>
        </Content>
      </Layout>
    </Layout>
  );
};

export default DashboardLayout;
