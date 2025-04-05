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

// Update interface definition to match new data format
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
  // Filter columns based on visible languages
  const filteredColumns = visibleLanguages
    ? columns.filter((col) => {
        // Keep the first column (key name) and the last column (action)
        if (col.key === "key_name" || col.key === "action") return true;
        // Display language columns based on user selection
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
          showTotal: (total) => `Total ${total} records`,
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

// Update column generation logic to adapt to new data format
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
      title: "Key Name",
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
                Has context
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
            Default
          </Tag>
        )}
      </div>
    ),
    dataIndex: ["languages", lang.code], // Access nested property
    key: lang.code,
    width: 200,
    render: (text: string, record: TranslationMatrix) => {
      // Get language value from new data structure
      const value = record.languages?.[lang.code];

      // If there's a translation value, show text; if not, show add button
      if (value) {
        // Find matching translation record (for edit/delete operations)
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
                    title="Are you sure you want to delete this translation?"
                    onConfirm={(e) => {
                      e?.stopPropagation();
                      onDeleteTranslation(translation.id);
                    }}
                    okText="Yes"
                    cancelText="Cancel"
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
            Add translation
          </Button>
        );
      }
    },
  }));

  // Action column
  const actionColumn = {
    title: "Action",
    key: "action",
    fixed: "right" as const,
    width: 120,
    render: (_: any, record: TranslationMatrix) => (
      <Space size="small">
        <Button
          type="primary"
          size="small"
          icon={<EditOutlined />}
          onClick={() => onShowBatchAddModal(record.key_name, record.context)}
        >
          Edit
        </Button>
      </Space>
    ),
  };

  // Combine all columns
  return [...baseColumns, ...languageColumns, actionColumn];
};

export default TranslationTable;
