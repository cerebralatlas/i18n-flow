import React from "react";
import { Select } from "antd";
import { GlobalOutlined } from "@ant-design/icons";
import i18n from "../i18n";
import { useLanguageStore } from "../stores/langugageStore";

const { Option } = Select;

const LanguageSelector: React.FC = () => {
  const { language, setLanguage } = useLanguageStore();

  const handleLanguageChange = (value: string) => {
    i18n.changeLanguage(value);
    localStorage.setItem("language", value);
    setLanguage(value);
  };

  return (
    <Select
      value={language}
      onChange={handleLanguageChange}
      style={{ width: 120 }}
      suffixIcon={<GlobalOutlined />}
    >
      <Option value="en">English</Option>
      <Option value="zh">中文</Option>
    </Select>
  );
};

export default LanguageSelector;
