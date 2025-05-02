import { Table, PaginationProps } from "antd";

import { TranslationResponse, Language } from "../../types/translation";
import { ColumnsType } from "antd/es/table";
import { AnyObject } from "antd/es/_util/type";
import { useTranslation } from "react-i18next";

interface TranslationLanguageInfo {
  id: number;
  value: string;
}

// Update interface definition to match new data format
export interface TranslationMatrix {
  key_name: string;
  context?: string;
  languages: {
    [languageCode: string]: TranslationLanguageInfo;
  };
}

interface TranslationTableProps {
  loading: boolean;
  paginatedMatrix: TranslationMatrix[];
  columns: ColumnsType<AnyObject>;
  translations: TranslationResponse[];
  languages: Language[];
  selectedKeys: string[];
  pagination: {
    current: number;
    pageSize: number;
    total: number;
  };
  onTableChange: (pagination: PaginationProps) => void;
  onRowSelectionChange: (
    selectedRowKeys: React.Key[],
    selectedRows: TranslationMatrix[]
  ) => void;
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
  const { t } = useTranslation();

  // 根据可见语言过滤列
  const filteredColumns = visibleLanguages
    ? columns.filter((col) => {
        // 保留第一列（key name）和最后一列（action）
        if (col.key === "key_name" || col.key === "action") return true;
        // 根据用户选择显示语言列
        return visibleLanguages.includes(col.key as string);
      })
    : columns;

  return (
    <div className="overflow-auto">
      <Table
        columns={filteredColumns as unknown as ColumnsType<TranslationMatrix>}
        dataSource={paginatedMatrix}
        rowKey="key_name"
        pagination={{
          ...pagination,
          showTotal: (total) => t("translation.table.total", { count: total }),
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

export default TranslationTable;
