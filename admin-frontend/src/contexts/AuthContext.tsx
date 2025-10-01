import React, { createContext, useState, useContext, useEffect } from "react";
import axios from "axios";
import { message } from "antd";
import { useTranslation } from "react-i18next";

interface User {
  id: number;
  username: string;
  email: string;
  role: 'admin' | 'member' | 'viewer';
  status: 'active' | 'disabled';
}

interface AuthContextType {
  user: User | null;
  loading: boolean;
  login: (username: string, password: string) => Promise<void>;
  logout: () => void;
  isAuthenticated: boolean;
  hasRole: (role: string) => boolean;
  isAdmin: () => boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);
  const { t } = useTranslation();

  // 检查用户认证状态
  const checkAuthStatus = async () => {
    try {
      setLoading(true);
      // 获取token
      const token = localStorage.getItem("authToken");

      if (!token) {
        setLoading(false);
        return;
      }

      // 设置axios默认请求头
      axios.defaults.headers.common["Authorization"] = `Bearer ${token}`;

      // 调用API验证当前用户
      const response = await axios.get("/api/user/info");
      setUser(response.data);
    } catch (error) {
      console.error("Validate user session failed:", error);
      // 清除无效token
      localStorage.removeItem("authToken");
      delete axios.defaults.headers.common["Authorization"];
    } finally {
      setLoading(false);
    }
  };

  // 组件挂载时检查认证状态
  useEffect(() => {
    checkAuthStatus();
  }, []);

  const login = async (username: string, password: string) => {
    try {
      const response = await axios.post("/api/login", {
        username,
        password,
      });
      const { token, user: userData } = response.data;

      // 保存token到localStorage
      localStorage.setItem("authToken", token);

      // 设置axios默认请求头
      axios.defaults.headers.common["Authorization"] = `Bearer ${token}`;

      setUser(userData);
    } catch (error) {
      console.error("Login failed:", error);
      throw error;
    }
  };

  const logout = () => {
    // 清除认证信息
    localStorage.removeItem("authToken");

    // 清除axios默认请求头
    delete axios.defaults.headers.common["Authorization"];

    message.success(t("common.messages.logoutSuccess"));
    setUser(null);
  };

  // 检查用户是否具有指定角色
  const hasRole = (role: string): boolean => {
    if (!user) return false;
    
    const roleLevel = {
      'viewer': 1,
      'member': 2,
      'admin': 3
    };
    
    const userLevel = roleLevel[user.role as keyof typeof roleLevel] || 0;
    const requiredLevel = roleLevel[role as keyof typeof roleLevel] || 0;
    
    return userLevel >= requiredLevel;
  };

  // 检查用户是否是管理员
  const isAdmin = (): boolean => {
    return user?.role === 'admin';
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        loading,
        login,
        logout,
        isAuthenticated: !!user,
        hasRole,
        isAdmin,
      }}
    >
      {children}
    </AuthContext.Provider>
  );
};

export const useAuth = () => {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error("useAuth must be used within an AuthProvider");
  }
  return context;
};
