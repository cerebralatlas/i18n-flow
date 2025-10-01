import {
  BrowserRouter as Router,
  Routes,
  Route,
  Navigate,
} from "react-router-dom";
import Login from "./pages/Login";
import DashboardLayout from "./pages/DashboardLayout";
import Dashboard from "./pages/Dashboard";
import ProjectManagement from "./pages/ProjectManagement";
import TranslationManagement from "./pages/TranslationManagement";
import UserManagement from "./pages/UserManagement";
import ProtectedRoute from "./components/ProtectedRoute";
import { useEffect } from "react";
import i18n from "./i18n";
import { useLanguageStore } from "./stores/langugageStore";
import { useAuthStore } from "./stores/authStore";

function App() {
  const { setLanguage } = useLanguageStore();
  const { checkAuthStatus } = useAuthStore();

  useEffect(() => {
    checkAuthStatus();
  }, [checkAuthStatus]);

  useEffect(() => {
    i18n.changeLanguage(localStorage.getItem("language") || i18n.language);
    setLanguage(localStorage.getItem("language") || i18n.language);
  }, [setLanguage]);

  return (
    <Router>
      <Routes>
        <Route path="/login" element={<Login />} />
        <Route
          path="/"
          element={
            <ProtectedRoute>
              <DashboardLayout />
            </ProtectedRoute>
          }
        >
          <Route index element={<Navigate to="dashboard" />} />
          <Route path="dashboard" element={<Dashboard />} />
          <Route path="projects" element={<ProjectManagement />} />
          <Route path="translations" element={<TranslationManagement />} />
          <Route
            path="translations/project/:projectId"
            element={<TranslationManagement />}
          />
          <Route path="users" element={<UserManagement />} />
        </Route>
      </Routes>
    </Router>
  );
}

export default App;
