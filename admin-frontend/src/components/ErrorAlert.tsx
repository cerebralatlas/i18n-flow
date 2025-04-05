import React from "react";
import { Alert, Space } from "antd";

interface ErrorAlertProps {
  message: string;
  description?: string;
  onRetry?: () => void;
}

const ErrorAlert: React.FC<ErrorAlertProps> = ({
  message,
  description = "Please try again later or contact the administrator",
  onRetry,
}) => {
  return (
    <Alert
      message={message}
      description={
        <Space direction="vertical">
          <span>{description}</span>
          {onRetry && <a onClick={onRetry}>Retry</a>}
        </Space>
      }
      type="error"
      showIcon
      closable
    />
  );
};

export default ErrorAlert;
