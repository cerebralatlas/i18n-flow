import React, { useState, useEffect, useCallback } from "react";
import {
  Button,
  Input,
  Select,
  Space,
  Upload,
  Popconfirm,
  message,
} from "antd";
import {
  PlusOutlined,
  SearchOutlined,
  DeleteOutlined,
  ImportOutlined,
  ExportOutlined,
  FileExcelOutlined,
} from "@ant-design/icons";
import { Project } from "../../types/project";
import { Language } from "../../types/translation";
import ColumnSelector from "./ColumnSelector";
import debounce from "lodash/debounce";
import { useTranslation } from "react-i18next";

interface TranslationToolbarProps {
  projects: Project[];
  selectedProject: number | null;
  keyword: string;
  onProjectChange: (projectId: number) => void;
  onKeywordChange: (keyword: string) => void;
  onSearch: () => void;
  onAddTranslation: () => void;
  onBatchAddTranslation: () => void;
  onImportJsonClick: () => void;
  onExportClick: () => void;
  onExcelFileUpload: (file: File) => boolean;
  selectedTranslations: number[];
  onBatchDelete: () => void;
  batchDeleteLoading: boolean;
  languages: Language[];
  selectedLanguageColumns: string[];
  onColumnSelectionChange: (selectedCodes: string[]) => void;
}

const TranslationToolbar: React.FC<TranslationToolbarProps> = ({
  projects,
  selectedProject,
  keyword,
  onProjectChange,
  onKeywordChange,
  onSearch,
  onAddTranslation,
  onBatchAddTranslation,
  onImportJsonClick,
  onExportClick,
  onExcelFileUpload,
  selectedTranslations,
  onBatchDelete,
  batchDeleteLoading,
  languages,
  selectedLanguageColumns,
  onColumnSelectionChange,
}) => {
  const [localKeyword, setLocalKeyword] = useState(keyword);
  const { t } = useTranslation();

  useEffect(() => {
    setLocalKeyword(keyword);
  }, [keyword]);

  // eslint-disable-next-line react-hooks/exhaustive-deps
  const debouncedSearch = useCallback(
    debounce((value: string) => {
      onKeywordChange(value);
      onSearch();
    }, 500),
    [onKeywordChange, onSearch]
  );

  const handleInputChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    const value = e.target.value;
    setLocalKeyword(value);
    debouncedSearch(value);
  };

  const handlePressEnter = () => {
    debouncedSearch.cancel();
    onKeywordChange(localKeyword);
    onSearch();
  };

  const handleProjectSelect = (value: number) => {
    onProjectChange(value);
  };

  return (
    <>
      <div className="mb-6 flex justify-between items-center">
        <h3 className="text-xl font-bold m-0">{t("translation.title")}</h3>
        <div>
          <Select
            style={{ width: 200, marginRight: 16 }}
            placeholder={t("translation.toolbar.selectProject")}
            value={selectedProject}
            onChange={handleProjectSelect}
            options={projects?.map((project) => ({
              value: project.id,
              label: project.name,
            }))}
          />
          <Input
            placeholder={t("translation.toolbar.searchPlaceholder")}
            value={localKeyword}
            onChange={handleInputChange}
            prefix={<SearchOutlined />}
            style={{ width: 250, marginRight: 16 }}
            allowClear
            onPressEnter={handlePressEnter}
            onBlur={handlePressEnter}
          />
        </div>
      </div>

      <div className="mb-4 flex flex-wrap gap-2 items-center justify-between">
        <div className="flex flex-wrap gap-2 items-center">
          {selectedTranslations.length > 0 && (
            <Popconfirm
              title={`${t("translation.toolbar.batchDelete")} ${
                selectedTranslations.length
              }?`}
              onConfirm={onBatchDelete}
              okText="Yes"
              cancelText="Cancel"
            >
              <Button
                type="primary"
                danger
                icon={<DeleteOutlined />}
                loading={batchDeleteLoading}
              >
                {t("translation.toolbar.batchDelete")}
              </Button>
            </Popconfirm>
          )}
        </div>
        <div className="flex flex-wrap gap-2 items-center">
          <ColumnSelector
            languages={languages}
            selectedColumns={selectedLanguageColumns}
            onChange={onColumnSelectionChange}
          />
          <Space size="small">
            <Button
              type="primary"
              icon={<PlusOutlined />}
              onClick={() => {
                if (!selectedProject) {
                  message.warning(t("translation.message.selectProject"));
                  return;
                }
                onAddTranslation();
              }}
            >
              {t("translation.toolbar.addTranslation")}
            </Button>

            <Button
              icon={<PlusOutlined />}
              onClick={() => {
                if (!selectedProject) {
                  message.warning(t("translation.message.selectProject"));
                  return;
                }
                onBatchAddTranslation();
              }}
            >
              {t("translation.toolbar.batchAdd")}
            </Button>

            <Button icon={<ImportOutlined />} onClick={onImportJsonClick}>
              {t("translation.toolbar.import.json")}
            </Button>

            <Upload
              beforeUpload={onExcelFileUpload}
              showUploadList={false}
              accept=".xlsx,.xls"
            >
              <Button icon={<FileExcelOutlined />}>
                {t("translation.toolbar.import.excel")}
              </Button>
            </Upload>

            <Button icon={<ExportOutlined />} onClick={onExportClick}>
              {t("translation.toolbar.export")}
            </Button>
          </Space>
        </div>
      </div>
    </>
  );
};

export default TranslationToolbar;
