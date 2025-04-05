import React, { useEffect } from "react";
import { Modal, Form, Input, Select, Upload, Table, Tag } from "antd";
import { UploadOutlined } from "@ant-design/icons";
import { Project } from "../../types/project";
import { Language, TranslationResponse } from "../../types/translation";
import { TranslationMatrix } from "./TranslationTable";

const { TextArea } = Input;
const { Dragger } = Upload;

interface BaseModalProps {
  visible: boolean;
  onCancel: () => void;
  projects: Project[];
  languages: Language[];
  selectedProject: number | null;
}

// Single Translation Modal
interface CreateModalProps extends BaseModalProps {
  form: any;
  onOk: () => void;
}

export const CreateTranslationModal: React.FC<CreateModalProps> = ({
  visible,
  onCancel,
  onOk,
  form,
  projects,
  languages,
  selectedProject,
}) => {
  return (
    <Modal
      title="添加翻译"
      open={visible}
      onOk={onOk}
      onCancel={onCancel}
      okText="创建"
      cancelText="取消"
      maskClosable={false}
    >
      <Form form={form} layout="vertical" name="translationForm">
        <Form.Item
          name="project_id"
          label="项目"
          rules={[{ required: true, message: "请选择项目" }]}
        >
          <Select
            placeholder="请选择项目"
            options={projects.map((project) => ({
              value: project.id,
              label: project.name,
            }))}
            disabled={!!selectedProject}
          />
        </Form.Item>

        <Form.Item
          name="key_name"
          label="键名"
          rules={[{ required: true, message: "请输入键名" }]}
        >
          <Input placeholder="请输入键名" />
        </Form.Item>

        <Form.Item name="context" label="上下文说明">
          <Input placeholder="请输入上下文说明（可选）" />
        </Form.Item>

        <Form.Item
          name="language_id"
          label="语言"
          rules={[{ required: true, message: "请选择语言" }]}
        >
          <Select
            placeholder="请选择语言"
            options={languages.map((lang) => ({
              value: lang.id,
              label: `${lang.name} [${lang.code}]`,
            }))}
          />
        </Form.Item>

        <Form.Item
          name="value"
          label="翻译值"
          rules={[{ required: true, message: "请输入翻译值" }]}
        >
          <TextArea rows={4} placeholder="请输入翻译值" />
        </Form.Item>
      </Form>
    </Modal>
  );
};

// Batch Add Modal
interface BatchModalProps extends BaseModalProps {
  form: any;
  onOk: () => void;
  translations: TranslationResponse[];
  paginatedMatrix: TranslationMatrix[];
}

export const BatchTranslationModal: React.FC<BatchModalProps> = ({
  visible,
  onCancel,
  onOk,
  form,
  languages,
  translations,
  paginatedMatrix,
}) => {
  // Add useEffect to ensure form values are set properly when the modal is shown
  useEffect(() => {
    if (visible && form) {
      // 添加一个小延迟，确保表单已经完全初始化
      setTimeout(() => {
        // Get the current key name from the form
        const keyName = form.getFieldValue("key_name");
        console.log("Setting values for key:", keyName);

        if (keyName && paginatedMatrix && paginatedMatrix.length > 0) {
          // 创建一个包含所有值的对象，然后一次性设置
          const formValues = { key_name: keyName };

          // Set form values for each language that has a translation
          languages.forEach((lang) => {
            if (paginatedMatrix) {
              const paginatedTranslation = paginatedMatrix.find(
                (t) => t.key_name === keyName && t.languages[lang.code]
              );
              if (paginatedTranslation) {
                const langKey = `lang_${lang.code}` as keyof typeof formValues;
                formValues[langKey] = paginatedTranslation.languages[lang.code];
                console.log(
                  `Setting lang_${lang.code}:`,
                  paginatedTranslation.languages[lang.code]
                );
              }
            }
          });

          // 使用setFieldsValue一次性设置所有值
          form.setFieldsValue(formValues);
        }
      }, 100); // 小延迟确保DOM已更新
    }
  }, [visible, form, paginatedMatrix, languages]);

  return (
    <Modal
      title="批量添加/更新翻译"
      open={visible}
      onOk={onOk}
      onCancel={onCancel}
      okText="保存"
      cancelText="取消"
      maskClosable={false}
      width={700}
      destroyOnClose={true}
    >
      <Form
        form={form}
        layout="vertical"
        name="batchTranslationForm"
        preserve={false}
      >
        <Form.Item
          name="key_name"
          label="键名"
          rules={[{ required: true, message: "请输入键名" }]}
          initialValue={form.getFieldValue("key_name")}
        >
          <Input placeholder="请输入键名" />
        </Form.Item>

        <Form.Item name="context" label="上下文说明">
          <Input placeholder="请输入上下文说明（可选）" />
        </Form.Item>

        <div className="bg-gray-50 p-3 mb-4 rounded">
          <h5 className="text-lg font-medium">各语言的翻译值</h5>
          <p className="text-gray-500">
            已有的翻译值已自动填充，您可以修改或补充其他语言的翻译
          </p>
        </div>

        {languages.map((lang) => {
          // When the modal is displayed, query the existing translation for the key name and language
          const keyName = form.getFieldValue("key_name");
          const existingTranslation = translations.find(
            (t) => t.key_name === keyName && t.language_code === lang.code
          );

          return (
            <Form.Item
              key={lang.id}
              name={`lang_${lang.code}`}
              label={
                <div className="flex items-center">
                  <span>{`${lang.name} [${lang.code}]${
                    lang.is_default ? " (默认语言)" : ""
                  }`}</span>
                  {existingTranslation && (
                    <Tag color="blue" className="ml-2">
                      已有翻译
                    </Tag>
                  )}
                </div>
              }
            >
              <TextArea
                rows={2}
                placeholder={`请输入 ${lang.name} 的翻译值${
                  existingTranslation ? "（已有翻译）" : "（可选）"
                }`}
              />
            </Form.Item>
          );
        })}
      </Form>
    </Modal>
  );
};

// Edit Modal
interface EditModalProps extends BaseModalProps {
  form: any;
  onOk: () => void;
}

export const EditTranslationModal: React.FC<EditModalProps> = ({
  visible,
  onCancel,
  onOk,
  form,
  languages,
}) => {
  return (
    <Modal
      title="编辑翻译"
      open={visible}
      onOk={onOk}
      onCancel={onCancel}
      okText="更新"
      cancelText="取消"
      maskClosable={false}
    >
      <Form form={form} layout="vertical" name="editTranslationForm">
        <Form.Item name="project_id" label="项目" hidden>
          <Input />
        </Form.Item>

        <Form.Item name="key_name" label="键名">
          <Input disabled />
        </Form.Item>

        <Form.Item name="language_id" label="语言">
          <Select
            disabled
            options={languages.map((lang) => ({
              value: lang.id,
              label: `${lang.name} [${lang.code}]`,
            }))}
          />
        </Form.Item>

        <Form.Item name="context" label="上下文说明">
          <Input placeholder="请输入上下文说明（可选）" />
        </Form.Item>

        <Form.Item
          name="value"
          label="翻译值"
          rules={[{ required: true, message: "请输入翻译值" }]}
        >
          <TextArea rows={4} placeholder="请输入翻译值" />
        </Form.Item>
      </Form>
    </Modal>
  );
};

// Import JSON Modal
interface ImportModalProps {
  visible: boolean;
  onCancel: () => void;
  onImport: (file: File) => boolean;
}

export const ImportTranslationModal: React.FC<ImportModalProps> = ({
  visible,
  onCancel,
  onImport,
}) => {
  return (
    <Modal title="导入翻译" open={visible} onCancel={onCancel} footer={null}>
      <div className="p-4">
        <p className="mb-4">请选择JSON格式的翻译文件导入，文件格式应为：</p>
        <pre className="bg-gray-100 p-3 rounded mb-4 text-sm">
          {`{
  "en": {
    "welcome": "Welcome",
    "hello": "Hello"
  },
  "zh-CN": {
    "welcome": "欢迎",
    "hello": "你好"
  }
}`}
        </pre>

        <Dragger
          name="file"
          multiple={false}
          accept=".json"
          beforeUpload={(file) => onImport(file)}
          showUploadList={false}
        >
          <p className="ant-upload-drag-icon">
            <UploadOutlined />
          </p>
          <p className="ant-upload-text">点击或拖拽文件到此区域上传</p>
          <p className="ant-upload-hint">支持单个JSON文件上传</p>
        </Dragger>
      </div>
    </Modal>
  );
};

// Excel Import Modal
interface ExcelModalProps {
  visible: boolean;
  onCancel: () => void;
  onOk: () => void;
  excelPreviewColumns: any[];
  excelPreviewData: any[];
  selectedLanguages: { [key: string]: string };
  onLanguageSelect: (columnKey: string, languageCode: string) => void;
  languages: Language[];
  loading: boolean;
}

export const ExcelImportModal: React.FC<ExcelModalProps> = ({
  visible,
  onCancel,
  onOk,
  excelPreviewColumns,
  excelPreviewData,
  selectedLanguages,
  onLanguageSelect,
  languages,
  loading,
}) => {
  return (
    <Modal
      title="Excel翻译导入"
      open={visible}
      onOk={onOk}
      onCancel={onCancel}
      width={900}
      okText="导入"
      cancelText="取消"
      confirmLoading={loading}
    >
      <div className="mb-4">
        <div className="mb-2 font-medium">Excel列与语言映射</div>
        <div className="text-gray-500 mb-4">
          请指定Excel中各列对应的语言，只有映射了语言的列才会被导入
        </div>

        <div className="grid grid-cols-3 gap-4">
          {excelPreviewColumns.slice(1).map((column) => (
            <div key={column.key} className="mb-2">
              <div className="mb-1">{column.title} 列对应:</div>
              <Select
                style={{ width: "100%" }}
                placeholder="选择语言"
                value={selectedLanguages[column.dataIndex]}
                onChange={(value) => onLanguageSelect(column.dataIndex, value)}
                allowClear
              >
                {languages.map((lang) => (
                  <Select.Option key={lang.id} value={lang.code}>
                    {lang.name} [{lang.code}]
                  </Select.Option>
                ))}
              </Select>
            </div>
          ))}
        </div>
      </div>

      <div className="mt-6">
        <div className="mb-2 font-medium">预览数据</div>
        <div className="text-gray-500 mb-4">
          显示前10行数据用于预览，实际导入会处理所有数据
        </div>

        <Table
          columns={excelPreviewColumns}
          dataSource={excelPreviewData.slice(0, 10)}
          bordered
          size="small"
          pagination={false}
          scroll={{ x: "max-content" }}
        />

        <div className="mt-2 text-gray-500">
          总计 {excelPreviewData.length} 行数据
        </div>
      </div>
    </Modal>
  );
};

// Export all modals
export default {
  CreateTranslationModal,
  BatchTranslationModal,
  EditTranslationModal,
  ImportTranslationModal,
  ExcelImportModal,
};
