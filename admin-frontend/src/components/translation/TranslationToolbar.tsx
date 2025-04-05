import React from "react";
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
import { TranslationResponse, Language } from "../../types/translation";
import ColumnSelector from "./ColumnSelector";

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
  selectedTranslations: TranslationResponse[];
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
  const handleProjectSelect = (value: number) => {
    onProjectChange(value);
  };

  return (
    <>
      <div className="mb-6 flex justify-between items-center">
        <h3 className="text-xl font-bold m-0">翻译管理</h3>
        <div>
          <Select
            style={{ width: 200, marginRight: 16 }}
            placeholder="请选择项目"
            value={selectedProject}
            onChange={handleProjectSelect}
            options={projects.map((project) => ({
              value: project.id,
              label: project.name,
            }))}
          />
          <Input
            placeholder="搜索键名或翻译值"
            value={keyword}
            onChange={(e) => onKeywordChange(e.target.value)}
            prefix={<SearchOutlined />}
            style={{ width: 250, marginRight: 16 }}
            allowClear
            onPressEnter={onSearch}
          />
        </div>
      </div>

      <div className="mb-4 flex flex-wrap gap-2 items-center justify-between">
        <div className="flex flex-wrap gap-2 items-center">
          {selectedTranslations.length > 0 && (
            <Popconfirm
              title={`确定要删除选中的 ${selectedTranslations.length} 条翻译吗？`}
              onConfirm={onBatchDelete}
              okText="确定"
              cancelText="取消"
            >
              <Button
                type="primary"
                danger
                icon={<DeleteOutlined />}
                loading={batchDeleteLoading}
              >
                批量删除
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
                  message.warning("请先选择项目");
                  return;
                }
                onAddTranslation();
              }}
            >
              添加翻译
            </Button>

            <Button
              icon={<PlusOutlined />}
              onClick={() => {
                if (!selectedProject) {
                  message.warning("请先选择项目");
                  return;
                }
                onBatchAddTranslation();
              }}
            >
              批量添加
            </Button>

            <Button icon={<ImportOutlined />} onClick={onImportJsonClick}>
              导入JSON
            </Button>

            <Upload
              beforeUpload={onExcelFileUpload}
              showUploadList={false}
              accept=".xlsx,.xls"
            >
              <Button icon={<FileExcelOutlined />}>导入Excel</Button>
            </Upload>

            <Button icon={<ExportOutlined />} onClick={onExportClick}>
              导出
            </Button>
          </Space>
        </div>
      </div>
    </>
  );
};

export default TranslationToolbar;
