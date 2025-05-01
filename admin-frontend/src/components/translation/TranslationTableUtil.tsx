import {
  EditOutlined,
  DeleteOutlined,
  PlusOutlined,
  InfoCircleOutlined,
} from "@ant-design/icons";
import { Language, TranslationResponse } from "../../types/translation";
import { TranslationMatrix } from "./TranslationTable";
import { Tooltip, Tag, Space, Button, Popconfirm, Typography } from "antd";

const { Text } = Typography;

export const generateTableColumns = (
  languages: Language[],
  translations: TranslationResponse[],
  onEditTranslation: (translation: TranslationResponse) => void,
  onDeleteTranslation: (id: number) => void,
  onAddTranslation: (keyName: string, languageId: number) => void,
  onShowBatchAddModal: (keyName: string, context?: string) => void,
  selectedLanguageCodes: string[] = []
) => {
  // 基础列：key name和context
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

  // 为每个语言创建一个列 - 如果提供了选中的语言，则过滤
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
      // 从新的数据结构中获取语言值
      const value = record.languages?.[lang.code];

      // 如果存在翻译值，则显示文本；否则显示添加按钮
      if (value) {
        // 查找匹配的翻译记录（用于编辑/删除操作）
        const translation = translations.find(
          (t) => t.key_name === record.key_name && t.language_code === lang.code
        );

        return (
          <div className="group relative">
            <div className="min-h-[40px] whitespace-pre-wrap mb-2">
              {value.value}
            </div>
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

  // 操作列
  const actionColumn = {
    title: "Action",
    key: "action",
    fixed: "right" as const,
    width: 120,
    render: (_, record: TranslationMatrix) => (
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

  // 合并所有列
  return [...baseColumns, ...languageColumns, actionColumn];
};
