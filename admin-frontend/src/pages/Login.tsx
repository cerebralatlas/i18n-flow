import React, { useState } from "react";
import { Form, Input, Button, Card, message, Divider, Checkbox } from "antd";
import { UserOutlined, LockOutlined, KeyOutlined } from "@ant-design/icons";
import { useNavigate } from "react-router-dom";
import { useAuth } from "../contexts/AuthContext";

const Login: React.FC = () => {
  const [loading, setLoading] = useState(false);
  const { login } = useAuth();
  const navigate = useNavigate();

  const onFinish = async (values: { username: string; password: string }) => {
    try {
      setLoading(true);
      await login(values.username, values.password);
      message.success("Login successful");
      navigate("/dashboard"); // redirect to dashboard after login
    } catch (error) {
      console.error("Login failed:", error);
      message.error("Invalid username or password, please try again");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen flex items-center justify-center bg-gradient-to-r from-blue-500 to-purple-600">
      <div className="w-full max-w-md px-4">
        <Card
          className="w-full shadow-2xl rounded-xl overflow-hidden border-0"
          bodyStyle={{ padding: "2rem" }}
        >
          <div className="text-center mb-8">
            <div className="inline-flex justify-center items-center mb-4 w-16 h-16 rounded-full bg-blue-100">
              <KeyOutlined className="text-3xl text-blue-600" />
            </div>
            <h1 className="text-2xl font-bold text-gray-800 mb-2">i18n-flow</h1>
            <p className="text-gray-500">Translation Management System</p>
          </div>

          <Divider className="mb-6">Secure Login</Divider>

          <Form
            name="login"
            initialValues={{ remember: true }}
            onFinish={onFinish}
            layout="vertical"
            size="large"
          >
            <Form.Item
              name="username"
              rules={[
                { required: true, message: "Please enter your username!" },
              ]}
            >
              <Input
                prefix={<UserOutlined className="text-gray-400" />}
                placeholder="Username"
                className="rounded-lg py-2"
              />
            </Form.Item>

            <Form.Item
              name="password"
              rules={[
                { required: true, message: "Please enter your password!" },
              ]}
            >
              <Input.Password
                prefix={<LockOutlined className="text-gray-400" />}
                placeholder="Password"
                className="rounded-lg py-2"
              />
            </Form.Item>

            <div className="flex justify-between items-center mb-4">
              <Form.Item name="remember" valuePropName="checked" noStyle>
                <Checkbox>Remember me</Checkbox>
              </Form.Item>
            </div>

            <Form.Item>
              <Button
                type="primary"
                htmlType="submit"
                className="w-full h-12 rounded-lg font-medium text-base bg-gradient-to-r from-blue-500 to-blue-700 border-0 shadow-md hover:shadow-lg transition-all"
                loading={loading}
              >
                {loading ? "Logging in..." : "Sign In"}
              </Button>
            </Form.Item>
          </Form>

          <div className="text-center text-gray-500 text-sm mt-8">
            <p>© {new Date().getFullYear()} i18n-flow. All rights reserved.</p>
          </div>
        </Card>
      </div>
    </div>
  );
};

export default Login;
