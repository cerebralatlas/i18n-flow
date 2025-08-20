import React, { useEffect } from "react";
import {
  Modal,
  Form,
  Input,
  Select,
  Upload,
  Table,
  Tag,
  FormInstance,
} from "antd";
import { UploadOutlined } from "@ant-design/icons";
import { Project } from "../../types/project";
import {
  ExcelData,
  ExcelPreviewColumns,
  Language,
  TranslationResponse,
} from "../../types/translation";
import { TranslationMatrix } from "./TranslationTable";
import { useTranslation } from "react-i18next";
import { ColumnsType } from "antd/es/table";

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
  form: FormInstance;
  onOk: () => void;
}

// 单个添加翻译的modal
export const CreateTranslationModal: React.FC<CreateModalProps> = ({
  visible,
  onCancel,
  onOk,
  form,
  projects,
  languages,
  selectedProject,
}) => {
  const { t } = useTranslation();
  return (
    <Modal
      title={t("translation.modal.create.title")}
      open={visible}
      onOk={onOk}
      onCancel={onCancel}
      okText={t("translation.modal.common.create")}
      cancelText={t("translation.modal.common.cancel")}
      maskClosable={false}
    >
      <Form form={form} layout="vertical" name="translationForm">
        <Form.Item
          name="project_id"
          label={t("translation.modal.common.project.label")}
          rules={[
            {
              required: true,
              message: t("translation.modal.common.project.required"),
            },
          ]}
        >
          <Select
            placeholder={t("translation.modal.common.project.placeholder")}
            options={projects?.map((project) => ({
              value: project.id,
              label: project.name,
            }))}
            disabled={!!selectedProject}
          />
        </Form.Item>

        <Form.Item
          name="key_name"
          label={t("translation.modal.common.keyName.label")}
          rules={[
            {
              required: true,
              message: t("translation.modal.common.keyName.required"),
            },
          ]}
        >
          <Input
            placeholder={t("translation.modal.common.keyName.placeholder")}
          />
        </Form.Item>

        <Form.Item
          name="context"
          label={t("translation.modal.common.context.label")}
        >
          <Input
            placeholder={t("translation.modal.common.context.placeholder")}
          />
        </Form.Item>

        <Form.Item
          name="language_id"
          label={t("translation.modal.common.language.label")}
          rules={[
            {
              required: true,
              message: t("translation.modal.common.language.required"),
            },
          ]}
        >
          <Select
            placeholder={t("translation.modal.common.language.placeholder")}
            options={languages.map((lang) => ({
              value: lang.id,
              label: t("translation.modal.common.language.format", {
                name: lang.name,
                code: lang.code,
              }),
            }))}
          />
        </Form.Item>

        <Form.Item
          name="value"
          label={t("translation.modal.common.value.label")}
          rules={[
            {
              required: true,
              message: t("translation.modal.common.value.required"),
            },
          ]}
        >
          <TextArea
            rows={4}
            placeholder={t("translation.modal.common.value.placeholder")}
          />
        </Form.Item>
      </Form>
    </Modal>
  );
};

// Batch Add Modal
interface BatchModalProps extends BaseModalProps {
  form: FormInstance;
  onOk: () => void;
  translations: TranslationResponse[];
  paginatedMatrix: TranslationMatrix[];
}

// 批量添加翻译的modal
export const BatchTranslationModal: React.FC<BatchModalProps> = ({
  visible,
  onCancel,
  onOk,
  form,
  languages,
  translations,
  paginatedMatrix,
}) => {
  const { t } = useTranslation();
  // 当modal显示时，确保form的值被正确设置
  useEffect(() => {
    if (visible && form) {
      // 添加一个小延迟，确保表单已经完全初始化
      setTimeout(() => {
        // 得到当前的key name
        const keyName = form.getFieldValue("key_name");

        if (keyName && paginatedMatrix && paginatedMatrix.length > 0) {
          // 创建一个包含所有值的对象，然后一次性设置
          const formValues = { key_name: keyName };

          // 遍历所有语言，设置form的值
          languages.forEach((lang) => {
            if (paginatedMatrix) {
              const paginatedTranslation = paginatedMatrix.find(
                (t) => t.key_name === keyName && t.languages[lang.code]
              );
              if (paginatedTranslation) {
                const langKey = `lang_${lang.code}` as keyof typeof formValues;
                formValues[langKey] =
                  paginatedTranslation.languages[lang.code].value;
              }
            }
          });

          // 使用setFieldsValue一次性设置所有值
          form.setFieldsValue(formValues);
        }
      }, 100);
    }
  }, [visible, form, paginatedMatrix, languages]);

  return (
    <Modal
      title={t("translation.modal.batch.title2")}
      open={visible}
      onOk={onOk}
      onCancel={onCancel}
      okText={t("translation.modal.common.save")}
      cancelText={t("translation.modal.common.cancel")}
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
          label={t("translation.modal.common.keyName.label")}
          rules={[
            {
              required: true,
              message: t("translation.modal.common.keyName.required"),
            },
          ]}
          initialValue={form.getFieldValue("key_name")}
        >
          <Input
            placeholder={t("translation.modal.common.keyName.placeholder")}
          />
        </Form.Item>

        <Form.Item
          name="context"
          label={t("translation.modal.common.context.label")}
        >
          <Input
            placeholder={t("translation.modal.common.context.placeholder")}
          />
        </Form.Item>

        <div className="bg-gray-50 p-3 mb-4 rounded">
          <h5 className="text-lg font-medium">
            {t("translation.modal.batch.valuesSection.title")}
          </h5>
          <p className="text-gray-500">
            {t("translation.modal.batch.valuesSection.description")}
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
                  <span>
                    {t("translation.modal.common.language.format", {
                      name: lang.name,
                      code: lang.code,
                    })}
                    {lang.is_default
                      ? t("translation.modal.batch.language.default")
                      : ""}
                  </span>
                  {existingTranslation && (
                    <Tag color="blue" className="ml-2">
                      {t("translation.modal.batch.language.existing")}
                    </Tag>
                  )}
                </div>
              }
            >
              <TextArea
                rows={2}
                placeholder={`${t(
                  "translation.modal.batch.language.placeholder",
                  { name: lang.name }
                )} ${
                  existingTranslation
                    ? t("translation.modal.batch.language.hasValue")
                    : t("translation.modal.batch.language.optional")
                }`}
              />
            </Form.Item>
          );
        })}
      </Form>
    </Modal>
  );
};

interface EditModalProps extends BaseModalProps {
  form: FormInstance;
  onOk: () => void;
}

// 编辑翻译的modal
export const EditTranslationModal: React.FC<EditModalProps> = ({
  visible,
  onCancel,
  onOk,
  form,
  languages,
}) => {
  const { t } = useTranslation();

  return (
    <Modal
      title={t("translation.modal.edit.title")}
      open={visible}
      onOk={onOk}
      onCancel={onCancel}
      okText={t("translation.modal.common.update")}
      cancelText={t("translation.modal.common.cancel")}
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

        <Form.Item
          name="context"
          label={t("translation.modal.edit.context.label")}
        >
          <Input
            placeholder={t("translation.modal.edit.context.placeholder")}
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

interface JSONImportModalProps {
  visible: boolean;
  onCancel: () => void;
  onImport: (file: File) => void;
}

// 导入json的modal
export const JSONImportTranslationModal: React.FC<JSONImportModalProps> = ({
  visible,
  onCancel,
  onImport,
}) => {
  const { t } = useTranslation();

  const handleUpload = (file: File) => {
    onImport(file);
    return false; // 阻止自动上传
  };

  return (
    <Modal
      title={t("translation.modal.import.jsonTitle")}
      open={visible}
      onCancel={onCancel}
      footer={null}
    >
      <div className="p-4">
        <p className="mb-4">{t("translation.modal.import.json.format")}</p>
        <pre className="bg-gray-100 p-3 rounded mb-4 text-sm">
          {`{
  "en": {
    "common.welcome": "Welcome",
    "common.hello": "Hello",
    "nav.home": "Home",
    "nav.about": "About"
  },
  "zh-CN": {
    "common.welcome": "欢迎",
    "common.hello": "你好",
    "nav.home": "首页",
    "nav.about": "关于"
  }
}`}
        </pre>

        <Dragger
          name="file"
          multiple={false}
          accept=".json"
          beforeUpload={handleUpload}
          showUploadList={false}
        >
          <p className="ant-upload-drag-icon">
            <UploadOutlined />
          </p>
          <p className="ant-upload-text">
            {t("translation.modal.import.json.uploadText")}
          </p>
          <p className="ant-upload-hint">
            {t("translation.modal.import.json.uploadHint")}
          </p>
        </Dragger>
      </div>
    </Modal>
  );
};

// 导入excel的modal
interface ExcelModalProps {
  visible: boolean;
  onCancel: () => void;
  onOk: () => void;
  excelPreviewColumns: ExcelPreviewColumns[];
  excelPreviewData: ExcelData[];
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
  const { t } = useTranslation();
  return (
    <Modal
      title={t("translation.modal.import.excelTitle")}
      open={visible}
      onOk={onOk}
      onCancel={onCancel}
      width={900}
      okText={t("translation.modal.common.save")}
      cancelText={t("translation.modal.common.cancel")}
      confirmLoading={loading}
    >
      <div className="mb-4">
        <div className="mb-2 font-medium">
          {t("translation.modal.import.excel.mapping.title")}
        </div>
        <div className="text-gray-500 mb-4">
          {t("translation.modal.import.excel.mapping.description")}
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
        <div className="mb-2 font-medium">
          {t("translation.modal.import.excel.preview.title")}
        </div>
        <div className="text-gray-500 mb-4">
          {t("translation.modal.import.excel.preview.description")}
        </div>

        <Table
          columns={excelPreviewColumns as unknown as ColumnsType<ExcelData>}
          dataSource={excelPreviewData.slice(0, 10)}
          bordered
          size="small"
          pagination={false}
          scroll={{ x: "max-content" }}
        />

        <div className="mt-2 text-gray-500">
          {t("translation.modal.import.excel.preview.total", {
            count: excelPreviewData.length,
          })}
        </div>
      </div>
    </Modal>
  );
};

// 导出所有modal
export default {
  CreateTranslationModal,
  BatchTranslationModal,
  EditTranslationModal,
  ImportTranslationModal: JSONImportTranslationModal,
  ExcelImportModal,
};
