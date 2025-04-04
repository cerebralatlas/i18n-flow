import {
  BrowserRouter as Router,
  Routes,
  Route,
  Navigate,
} from "react-router-dom";
import { AuthProvider } from "./contexts/AuthContext";
import Login from "./pages/Login";
import DashboardLayout from "./pages/DashboardLayout";
import Dashboard from "./pages/Dashboard";
import ProjectManagement from "./pages/ProjectManagement";
import TranslationManagement from "./pages/TranslationManagement";
import ProtectedRoute from "./components/ProtectedRoute";

function App() {
  return (
    <AuthProvider>
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
            {/* 嵌套路由 */}
            <Route index element={<Navigate to="dashboard" />} />
            <Route path="dashboard" element={<Dashboard />} />
            <Route path="projects" element={<ProjectManagement />} />
            <Route path="translations" element={<TranslationManagement />} />
            <Route
              path="translations/project/:projectId"
              element={<TranslationManagement />}
            />
            {/* 其他路由可以在这里添加 */}
          </Route>
        </Routes>
      </Router>
    </AuthProvider>
  );
}

export default App;
