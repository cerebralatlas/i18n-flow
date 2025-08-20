import React from "react";
import { Alert, Space } from "antd";
import { useTranslation } from "react-i18next";

interface ErrorAlertProps {
  message: string;
  description?: string;
  onRetry?: () => void;
}

const ErrorAlert: React.FC<ErrorAlertProps> = ({
  message,
  description,
  onRetry,
}) => {
  const { t } = useTranslation();
  const defaultDescription = description || t("common.error.defaultDescription");
  return (
    <Alert
      message={message}
      description={
        <Space direction="vertical">
          <span>{defaultDescription}</span>
          {onRetry && <a onClick={onRetry}>{t("common.error.retry")}</a>}
        </Space>
      }
      type="error"
      showIcon
      closable
    />
  );
};

export default ErrorAlert;
