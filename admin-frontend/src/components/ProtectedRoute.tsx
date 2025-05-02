import React from "react";
import { Navigate } from "react-router-dom";
import { useAuthStore } from "../stores/authStore";
import { Spin } from "antd";

interface ProtectedRouteProps {
  children: React.ReactNode;
}

const ProtectedRoute: React.FC<ProtectedRouteProps> = ({ children }) => {
  const { isAuthenticated, loading } = useAuthStore();

  // If authentication status is being checked, show loading indicator
  if (loading) {
    return (
      <div className="min-h-screen flex justify-center items-center">
        <Spin size="large" tip="Verifying identity..." />
      </div>
    );
  }

  // After checking authentication status, if not authenticated, redirect to login page
  if (!isAuthenticated) {
    return <Navigate to="/login" replace />;
  }

  // Authenticated, render the child component
  return <>{children}</>;
};

export default ProtectedRoute;
