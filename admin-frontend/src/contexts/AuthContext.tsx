import React, { createContext, useState, useContext, useEffect } from "react";
import axios from "axios";

interface User {
  id: number;
  username: string;
  // Other user attributes...
}

interface AuthContextType {
  user: User | null;
  loading: boolean;
  login: (username: string, password: string) => Promise<void>;
  logout: () => void;
  isAuthenticated: boolean;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export const AuthProvider: React.FC<{ children: React.ReactNode }> = ({
  children,
}) => {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true);

  // Add a method to check user authentication status
  const checkAuthStatus = async () => {
    try {
      setLoading(true);
      // Get token from localStorage
      const token = localStorage.getItem("authToken");

      if (!token) {
        setLoading(false);
        return;
      }

      // Set axios default request header with token
      axios.defaults.headers.common["Authorization"] = `Bearer ${token}`;

      // Call API to validate current user
      const response = await axios.get("/api/user/info");
      setUser(response.data);
    } catch (error) {
      console.error("Validate user session failed:", error);
      // Clear invalid token
      localStorage.removeItem("authToken");
      delete axios.defaults.headers.common["Authorization"];
    } finally {
      setLoading(false);
    }
  };

  // Check authentication status when component mounts
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

      // Save token to localStorage
      localStorage.setItem("authToken", token);

      // Set axios default request header
      axios.defaults.headers.common["Authorization"] = `Bearer ${token}`;

      setUser(userData);
    } catch (error) {
      console.error("Login failed:", error);
      throw error;
    }
  };

  const logout = () => {
    // Clear authentication information in localStorage
    localStorage.removeItem("authToken");

    // Clear axios default request header
    delete axios.defaults.headers.common["Authorization"];

    setUser(null);
  };

  return (
    <AuthContext.Provider
      value={{
        user,
        loading,
        login,
        logout,
        isAuthenticated: !!user,
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
