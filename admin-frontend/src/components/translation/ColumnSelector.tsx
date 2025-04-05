import React from "react";
import { Button, Dropdown, Checkbox, Space, Typography } from "antd";
import { ColumnHeightOutlined } from "@ant-design/icons";
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
              <Text strong>显示列</Text>
              <Space>
                <Button size="small" onClick={handleSelectAll}>
                  全选
                </Button>
                <Button size="small" onClick={handleSelectNone}>
                  清除
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
                      {lang.is_default && <Text type="success">(默认)</Text>}
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
        显示列 ({selectedColumns.length}/{languages.length})
      </Button>
    </Dropdown>
  );
};

export default ColumnSelector;
