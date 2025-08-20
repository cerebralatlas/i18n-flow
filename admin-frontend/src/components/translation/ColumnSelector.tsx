import React from "react";
import { Button, Dropdown, Checkbox, Space, Typography } from "antd";
import { ColumnHeightOutlined } from "@ant-design/icons";
import { useTranslation } from "react-i18next";
import { Language } from "../../types/translation";

const { Text } = Typography;

interface ColumnSelectorProps {
  languages: Language[];
  selectedColumns: string[];
  onChange: (selectedCodes: string[]) => void;
}

const ColumnSelector: React.FC<ColumnSelectorProps> = ({
  languages,
  selectedColumns,
  onChange,
}) => {
  const { t } = useTranslation();
  const handleCheckboxChange = (code: string, checked: boolean) => {
    const newSelectedColumns = checked
      ? [...selectedColumns, code]
      : selectedColumns.filter((col) => col !== code);
    onChange(newSelectedColumns);
  };

  const handleSelectAll = () => {
    onChange(languages.map((lang) => lang.code));
  };

  const handleSelectNone = () => {
    onChange([]);
  };

  const menu = {
    items: [
      {
        key: "column-options",
        label: (
          <div className="p-2" onClick={(e) => e.stopPropagation()}>
            <div className="mb-2 flex justify-between">
              <Text strong>{t("translation.columnSelector.showColumns")}</Text>
              <Space>
                <Button size="small" onClick={handleSelectAll}>
                  {t("translation.columnSelector.selectAll")}
                </Button>
                <Button size="small" onClick={handleSelectNone}>
                  {t("translation.columnSelector.clear")}
                </Button>
              </Space>
            </div>
            <div className="max-h-[300px] overflow-y-auto">
              {languages.map((lang) => (
                <div key={lang.code} className="mb-2">
                  <Checkbox
                    checked={selectedColumns.includes(lang.code)}
                    onChange={(e) =>
                      handleCheckboxChange(lang.code, e.target.checked)
                    }
                  >
                    <Space>
                      {lang.name}
                      <Text type="secondary">[{lang.code}]</Text>
                      {lang.is_default && <Text type="success">({t("translation.columnSelector.default")})</Text>}
                    </Space>
                  </Checkbox>
                </div>
              ))}
            </div>
          </div>
        ),
      },
    ],
  };

  return (
    <Dropdown menu={menu} trigger={["click"]} placement="bottomRight">
      <Button icon={<ColumnHeightOutlined />}>
        {t("translation.columnSelector.buttonText", {
          selected: selectedColumns.length,
          total: languages.length
        })}
      </Button>
    </Dropdown>
  );
};

export default ColumnSelector;
