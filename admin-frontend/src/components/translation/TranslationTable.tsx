import {
  Table,
  Button,
  Space,
  Popconfirm,
  Tag,
  Tooltip,
  Typography,
} from "antd";
import {
  EditOutlined,
  DeleteOutlined,
  PlusOutlined,
  InfoCircleOutlined,
} from "@ant-design/icons";
import { TranslationResponse, Language } from "../../types/translation";

const { Text } = Typography;

// 更新接口定义以匹配新的数据格式
export interface TranslationMatrix {
  key_name: string;
  context?: string;
  languages: {
    [languageCode: string]: string;
  };
}

interface TranslationTableProps {
  loading: boolean;
  paginatedMatrix: TranslationMatrix[];
  columns: any[];
  translations: TranslationResponse[];
  languages: Language[];
  selectedKeys: string[];
  pagination: {
    current: number;
    pageSize: number;
    total: number;
  };
  onTableChange: (pagination: any) => void;
  onRowSelectionChange: (
    selectedRowKeys: React.Key[],
    selectedRows: TranslationMatrix[]
  ) => void;
  onEditTranslation: (translation: TranslationResponse) => void;
  onDeleteTranslation: (id: number) => void;
  onAddTranslation: (keyName: string, languageId: number) => void;
  onShowBatchAddModal: (keyName: string, context?: string) => void;
  visibleLanguages?: string[];
}

const TranslationTable: React.FC<TranslationTableProps> = ({
  loading,
  paginatedMatrix,
  columns,
  pagination,
  selectedKeys,
  onTableChange,
  onRowSelectionChange,
  visibleLanguages,
}) => {
  // 根据可见语言过滤列
  const filteredColumns = visibleLanguages
    ? columns.filter((col) => {
        // 保留第一列(键名)和最后一列(操作)
        if (col.key === "key_name" || col.key === "action") return true;
        // 根据用户选择显示语言列
        return visibleLanguages.includes(col.key);
      })
    : columns;

  return (
    <div className="overflow-auto">
      <Table
        columns={filteredColumns}
        dataSource={paginatedMatrix}
        rowKey="key_name"
        pagination={{
          ...pagination,
          showTotal: (total) => `共 ${total} 条记录`,
        }}
        onChange={onTableChange}
        scroll={{ x: "max-content" }}
        bordered
        size="middle"
        loading={loading}
        rowSelection={{
          type: "checkbox",
          selectedRowKeys: selectedKeys,
          onChange: onRowSelectionChange,
        }}
      />
    </div>
  );
};

// 更新列生成逻辑以适配新的数据格式
export const generateTableColumns = (
  languages: Language[],
  translations: TranslationResponse[],
  onEditTranslation: (translation: TranslationResponse) => void,
  onDeleteTranslation: (id: number) => void,
  onAddTranslation: (keyName: string, languageId: number) => void,
  onShowBatchAddModal: (keyName: string, context?: string) => void,
  selectedLanguageCodes: string[] = []
) => {
  // Base columns: key name and context
  const baseColumns = [
    {
      title: "键名",
      dataIndex: "key_name",
      key: "key_name",
      fixed: "left" as const,
      width: 200,
      render: (text: string, record: TranslationMatrix) => (
        <Space direction="vertical" size="small">
          <Text strong copyable>
            {text}
          </Text>
          {record.context && (
            <Tooltip title={record.context}>
              <Tag icon={<InfoCircleOutlined />} color="blue">
                有上下文说明
              </Tag>
            </Tooltip>
          )}
        </Space>
      ),
    },
  ];

  // Create a column for each language - filter by selected languages if provided
  const filteredLanguages =
    selectedLanguageCodes.length > 0
      ? languages.filter((lang) => selectedLanguageCodes.includes(lang.code))
      : languages;

  const languageColumns = filteredLanguages.map((lang) => ({
    title: (
      <div>
        <div>{lang.name}</div>
        <div className="text-xs text-gray-500">[{lang.code}]</div>
        {lang.is_default && (
          <Tag color="green" className="mt-1">
            默认
          </Tag>
        )}
      </div>
    ),
    dataIndex: ["languages", lang.code], // 使用嵌套属性访问
    key: lang.code,
    width: 200,
    render: (text: string, record: TranslationMatrix) => {
      // 从新的数据结构中获取语言值
      const value = record.languages?.[lang.code];

      // If there's a translation value, show text; if not, show add button
      if (value) {
        // 查找匹配的翻译记录（用于编辑/删除操作）
        const translation = translations.find(
          (t) => t.key_name === record.key_name && t.language_code === lang.code
        );

        return (
          <div className="group relative">
            <div className="min-h-[40px] whitespace-pre-wrap mb-2">{value}</div>
            {translation && (
              <div className="absolute top-0 right-0 opacity-0 group-hover:opacity-100 transition-opacity">
                <Space>
                  <Button
                    type="text"
                    size="small"
                    icon={<EditOutlined />}
                    onClick={(e) => {
                      e.stopPropagation();
                      onEditTranslation(translation);
                    }}
                  />
                  <Popconfirm
                    title="确定要删除该翻译吗？"
                    onConfirm={(e) => {
                      e?.stopPropagation();
                      onDeleteTranslation(translation.id);
                    }}
                    okText="确定"
                    cancelText="取消"
                  >
                    <Button
                      type="text"
                      size="small"
                      danger
                      icon={<DeleteOutlined />}
                      onClick={(e) => e.stopPropagation()}
                    />
                  </Popconfirm>
                </Space>
              </div>
            )}
          </div>
        );
      } else {
        return (
          <Button
            type="dashed"
            size="small"
            icon={<PlusOutlined />}
            onClick={() => onAddTranslation(record.key_name, lang.id)}
          >
            添加翻译
          </Button>
        );
      }
    },
  }));

  // Action column
  const actionColumn = {
    title: "操作",
    key: "action",
    fixed: "right" as const,
    width: 120,
    render: (_: any, record: TranslationMatrix) => (
      <Space size="small">
        <Button
          type="primary"
          size="small"
          icon={<PlusOutlined />}
          onClick={() => onShowBatchAddModal(record.key_name, record.context)}
        >
          批量添加
        </Button>
      </Space>
    ),
  };

  // Combine all columns
  return [...baseColumns, ...languageColumns, actionColumn];
};

export default TranslationTable;
