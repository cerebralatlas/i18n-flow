import React, { useState, useEffect } from "react";
import { Empty, Form, message, PaginationProps, Spin } from "antd";
import { useParams } from "react-router-dom";
import { useTranslation } from "react-i18next";

// Custom hooks and components
import { useTranslationData } from "../hooks/useTranslationData";
import TranslationTable from "../components/translation/TranslationTable";
import TranslationToolbar from "../components/translation/TranslationToolbar";
import {
  CreateTranslationModal,
  BatchTranslationModal,
  JSONImportTranslationModal,
  ExcelImportModal,
} from "../components/translation/TranslationModal";
import {
  ExcelData as ExcelDataUtil,
  parseExcelFile,
  formatExcelDataForImport,
  autoMapLanguageColumns,
} from "../components/translation/ExcelUtils";
import {
  loadSelectedColumns,
  saveSelectedColumns,
} from "../utils/localStorage";
import { processTranslationJSON } from "../utils/jsonFlattener";

// Types
import {
  BatchTranslationRequest,
  ExcelPreviewColumns,
  ExcelData as ExcelDataPreview,
} from "../types/translation";
import { generateTableColumns } from "../components/translation/TranslationTableUtil";
import { ColumnProps, ColumnsType } from "antd/es/table";
import { AnyObject } from "antd/es/_util/type";

const TranslationManagement: React.FC = () => {
  const { projectId } = useParams<{ projectId: string }>();
  const { t } = useTranslation();

  // Use the custom hook to handle translation data
  const {
    projects,
    languages,
    columns,
    selectedProject,
    loading,
    keyword,
    paginatedMatrix,
    localPagination,
    selectedKeys,
    batchDeleteLoading,
    selectedTranslations,

    setColumns,
    setSelectedProject,
    setKeyword,
    handleTableChange,
    fetchTranslations,
    createTranslation,
    batchCreateTranslations,
    deleteTranslation,
    exportTranslations,
    importTranslationsFromJson,
    batchDeleteTranslations,
    handleRowSelectionChange,
  } = useTranslationData(projectId);

  // 选择语言列
  const [selectedLanguageColumns, setSelectedLanguageColumns] = useState<
    string[]
  >([]);

  // 创建翻译模态框
  const [createModalVisible, setCreateModalVisible] = useState<boolean>(false);
  // 批量翻译模态框
  const [batchModalVisible, setBatchModalVisible] = useState<boolean>(false);
  // 导入翻译模态框
  const [importModalVisible, setImportModalVisible] = useState<boolean>(false);

  // excel导入模态框
  const [excelImportModalVisible, setExcelImportModalVisible] =
    useState<boolean>(false);
  // excel导入数据
  const [excelData, setExcelData] = useState<ExcelDataUtil[]>([]);
  // excel导入预览列
  const [excelPreviewColumns, setExcelPreviewColumns] = useState<
    ExcelPreviewColumns[]
  >([]);
  // excel导入预览数据
  const [excelPreviewData, setExcelPreviewData] = useState<ExcelDataPreview[]>(
    []
  );
  // 选择语言
  const [selectedLanguages, setSelectedLanguages] = useState<{
    [key: string]: string;
  }>({});
  // excel导入loading
  const [excelImportLoading, setExcelImportLoading] = useState<boolean>(false);

  // Forms
  const [singleForm] = Form.useForm();
  const [batchForm] = Form.useForm();

  // 可见语言
  const [visibleLanguages, setVisibleLanguages] = useState<string[]>([]);

  // Effect to initialize selected columns when languages load
  useEffect(() => {
    if (languages.length > 0) {
      // Try to load saved preferences from localStorage first
      const savedColumns = loadSelectedColumns();

      if (savedColumns) {
        // Filter saved columns to only include languages that still exist
        const validSavedColumns = savedColumns.filter((code) =>
          languages.some((lang) => lang.code === code)
        );

        if (validSavedColumns.length > 0) {
          setSelectedLanguageColumns(validSavedColumns);
        } else {
          // If no valid saved columns, select all languages
          setSelectedLanguageColumns(languages.map((lang) => lang.code));
        }
      } else {
        // If no saved columns, select all languages
        setSelectedLanguageColumns(languages.map((lang) => lang.code));
      }
    }
  }, [languages]);

  // Effect to生成表格列
  useEffect(() => {
    if (languages.length > 0) {
      const generatedColumns = generateTableColumns(
        languages,
        [],
        () => {},
        handleDeleteTranslation,
        handleAddTranslation,
        showBatchAddModal,
        selectedLanguageColumns // 传递选中的列到过滤
      );
      setColumns(generatedColumns as unknown as ColumnProps<AnyObject>[]);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [languages, selectedLanguageColumns]);

  // 初始化可见语言
  useEffect(() => {
    if (languages.length > 0 && visibleLanguages.length === 0) {
      setVisibleLanguages(languages.map((lang) => lang.code));
    }
  }, [languages, visibleLanguages.length]);

  // 列选择
  const handleColumnSelectionChange = (selectedCodes: string[]) => {
    setSelectedLanguageColumns(selectedCodes);
    saveSelectedColumns(selectedCodes);
  };

  // 增加翻译
  const handleAddTranslation = (keyName: string, languageId: number) => {
    if (!selectedProject) {
      message.warning(t("translation.message.selectProject"));
      return;
    }

    // 获取上下文信息
    const context = paginatedMatrix.find(
      (m) => m.key_name === keyName
    )?.context;

    singleForm.setFieldsValue({
      project_id: selectedProject,
      key_name: keyName,
      context: context,
      language_id: languageId,
    });
    setCreateModalVisible(true);
  };

  const handleCreateTranslation = async () => {
    try {
      const values = await singleForm.validateFields();
      const success = await createTranslation(values);
      if (success) {
        setCreateModalVisible(false);
        singleForm.resetFields();
      }
    } catch (error) {
      console.error("Create translation failed:", error);
    }
  };

  // 展示批量翻译模态框
  const showBatchAddModal = async (keyName: string, context?: string) => {
    if (!selectedProject) {
      message.warning(t("translation.message.selectProject"));
      return;
    }

    try {
      batchForm.resetFields();

      batchForm.setFieldsValue({
        project_id: selectedProject,
        key_name: keyName,
        context: context,
      });

      setBatchModalVisible(true);
    } catch (error) {
      console.error("Load existing translations failed:", error);
      message.error("Load existing translations failed");
    }
  };

  // 创建批量翻译
  const handleBatchCreateTranslations = async () => {
    try {
      const values = await batchForm.validateFields();

      const request: BatchTranslationRequest = {
        project_id: selectedProject!,
        key_name: values.key_name,
        context: values.context,
        translations: {},
      };

      // 处理表单中的语言字段
      Object.keys(values).forEach((key) => {
        if (key.startsWith("lang_")) {
          const langCode = key.replace("lang_", "");
          request.translations[langCode] = values[key] || "";
        }
      });

      // 只有当有值时才发送请求
      if (Object.keys(request.translations).length > 0) {
        const success = await batchCreateTranslations(request);
        if (success) {
          setBatchModalVisible(false);
          batchForm.resetFields();
        }
      } else {
        message.warning(t("translation.message.addLanguage"));
      }
    } catch (error) {
      console.error("Batch create translations failed:", error);
    }
  };

  // 删除翻译
  const handleDeleteTranslation = async (id: number) => {
    await deleteTranslation(id);
  };

  // 导出翻译
  const handleExportTranslations = async (format: string = "json") => {
    const data = await exportTranslations(format);
    if (data) {
      // 创建和下载文件
      const fileName = `translations_${selectedProject}_${new Date().toISOString()}.json`;
      const jsonStr = JSON.stringify(data, null, 2);
      const blob = new Blob([jsonStr], { type: "application/json" });
      const href = URL.createObjectURL(blob);

      const link = document.createElement("a");
      link.href = href;
      link.download = fileName;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
      URL.revokeObjectURL(href);
    }
  };

  // 导入翻译
  const handleImportTranslations = async (file: File) => {
    if (!selectedProject) {
      message.warning(t("translation.message.selectProject"));
      return;
    }

    try {
      const reader = new FileReader();
      reader.onload = async (e) => {
        try {
          const content = e.target?.result as string;
          const jsonData = JSON.parse(content);

          // 使用 processTranslationJSON 处理嵌套的 JSON 数据
          const processedData = processTranslationJSON(jsonData);

          const success = await importTranslationsFromJson(processedData);
          if (success) {
            setImportModalVisible(false);
            message.success(t("translation.message.importSuccess"));
            // 刷新翻译列表
            await fetchTranslations();
          }
        } catch (error) {
          console.error("Import failed:", error);
          message.error(
            t("translation.message.importFailed", {
              error:
                error instanceof Error ? error.message : "Invalid JSON format",
            })
          );
        }
      };
      reader.readAsText(file);
    } catch (error) {
      console.error("Import translations failed:", error);
      message.error(
        t("translation.message.importFailed", {
          error: "File reading failed",
        })
      );
    }
  };

  // Excel处理函数
  const handleExcelFile = (file: File) => {
    parseExcelFile(
      file,
      (data: ExcelDataUtil | null) => {
        if (data) {
          setExcelData(data.jsonData);
          setExcelPreviewColumns(data.columns);
          setExcelPreviewData(
            data.previewData as unknown as ExcelDataPreview[]
          );

          // 自动映射列到语言
          const initialMapping = autoMapLanguageColumns(
            data.columns,
            languages
          );
          setSelectedLanguages(initialMapping);

          setExcelImportModalVisible(true);
        }
      },
      setExcelImportLoading
    );
    return false; // 阻止默认上传行为
  };

  // 选择语言
  const handleLanguageSelect = (columnKey: string, languageCode: string) => {
    setSelectedLanguages((prev) => ({
      ...prev,
      [columnKey]: languageCode,
    }));
  };

  // excel导入
  const handleExcelImport = async () => {
    if (!selectedProject) {
      message.warning(t("translation.message.selectProject"));
      return;
    }

    if (Object.keys(selectedLanguages).length === 0) {
      message.warning(t("translation.message.addLanguage"));
      return;
    }

    try {
      setExcelImportLoading(true);

      // Format Excel data for import
      const importData = formatExcelDataForImport(excelData, selectedLanguages);

      // 使用 processTranslationJSON 处理数据
      const processedData = processTranslationJSON(importData);

      // 调用导入 API
      const success = await importTranslationsFromJson(processedData);
      if (success) {
        setExcelImportModalVisible(false);
        message.success(t("translation.message.importSuccess"));
        // 刷新翻译列表
        await fetchTranslations();
      }
    } catch (error) {
      console.error("Import Excel translations failed:", error);
      message.error(
        t("translation.message.importFailed", {
          error: error instanceof Error ? error.message : "Unknown error",
        })
      );
    } finally {
      setExcelImportLoading(false);
    }
  };

  return (
    <div className="bg-white p-6 rounded-lg shadow">
      {/* 翻译工具按钮 */}
      <TranslationToolbar
        projects={projects}
        selectedProject={selectedProject}
        keyword={keyword}
        onProjectChange={setSelectedProject}
        onKeywordChange={setKeyword}
        onSearch={fetchTranslations}
        onAddTranslation={() => {
          if (!selectedProject) {
            message.warning(t("translation.message.selectProject"));
            return;
          }
          singleForm.setFieldsValue({ project_id: selectedProject });
          setCreateModalVisible(true);
        }}
        onBatchAddTranslation={() => {
          if (!selectedProject) {
            message.warning(t("translation.message.selectProject"));
            return;
          }
          batchForm.setFieldsValue({ project_id: selectedProject });
          setBatchModalVisible(true);
        }}
        onImportJsonClick={() => setImportModalVisible(true)}
        onExportClick={() => handleExportTranslations()}
        onExcelFileUpload={handleExcelFile}
        selectedTranslations={selectedTranslations}
        onBatchDelete={batchDeleteTranslations}
        batchDeleteLoading={batchDeleteLoading}
        languages={languages}
        selectedLanguageColumns={selectedLanguageColumns}
        onColumnSelectionChange={handleColumnSelectionChange}
      />

      {/* 表格 */}
      {loading ? (
        <div className="flex justify-center items-center py-12">
          <Spin size="large" tip={t("translation.table.loading")} />
        </div>
      ) : paginatedMatrix.length > 0 ? (
        <TranslationTable
          loading={loading}
          paginatedMatrix={paginatedMatrix}
          columns={columns as unknown as ColumnsType<AnyObject>}
          translations={[]}
          languages={languages}
          selectedKeys={selectedKeys}
          pagination={localPagination}
          onTableChange={
            handleTableChange as (pagination: PaginationProps) => void
          }
          onRowSelectionChange={handleRowSelectionChange}
          onDeleteTranslation={handleDeleteTranslation}
          onAddTranslation={handleAddTranslation}
          onShowBatchAddModal={showBatchAddModal}
        />
      ) : (
        <Empty description={t("translation.table.noData")} />
      )}

      {/* 创建翻译弹窗 */}
      <CreateTranslationModal
        visible={createModalVisible}
        onCancel={() => setCreateModalVisible(false)}
        onOk={handleCreateTranslation}
        form={singleForm}
        projects={projects}
        languages={languages}
        selectedProject={selectedProject}
      />

      {/* 批量创建翻译弹窗 */}
      <BatchTranslationModal
        visible={batchModalVisible}
        onCancel={() => setBatchModalVisible(false)}
        onOk={handleBatchCreateTranslations}
        form={batchForm}
        projects={projects}
        languages={languages}
        selectedProject={selectedProject}
        translations={[]}
        paginatedMatrix={paginatedMatrix}
      />

      {/* json导入翻译弹窗 */}
      <JSONImportTranslationModal
        visible={importModalVisible}
        onCancel={() => setImportModalVisible(false)}
        onImport={handleImportTranslations}
      />

      {/* excel导入弹窗 */}
      <ExcelImportModal
        visible={excelImportModalVisible}
        onCancel={() => setExcelImportModalVisible(false)}
        onOk={handleExcelImport}
        excelPreviewColumns={excelPreviewColumns}
        excelPreviewData={excelPreviewData as ExcelDataPreview[]}
        selectedLanguages={selectedLanguages}
        onLanguageSelect={handleLanguageSelect}
        languages={languages}
        loading={excelImportLoading}
      />
    </div>
  );
};

export default TranslationManagement;
