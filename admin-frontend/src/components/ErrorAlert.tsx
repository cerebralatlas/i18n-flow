import React from "react";
import { Alert, Space } from "antd";

interface ErrorAlertProps {
  message: string;
  description?: string;
  onRetry?: () => void;
}

const ErrorAlert: React.FC<ErrorAlertProps> = ({
  message,
  description = "请稍后重试或联系管理员",
  onRetry,
}) => {
  return (
    <Alert
      message={message}
      description={
        <Space direction="vertical">
          <span>{description}</span>
          {onRetry && <a onClick={onRetry}>重试</a>}
        </Space>
      }
      type="error"
      showIcon
      closable
    />
  );
};

export default ErrorAlert;
