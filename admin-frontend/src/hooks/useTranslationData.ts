import { useState, useEffect } from "react";
import { message } from "antd";
import { translationService } from "../services/translationService";
import { projectService } from "../services/projectService";
import {
  Language,
  BatchTranslationRequest,
  TranslationTablePagination,
  TranslationRequest,
} from "../types/translation";
import { Project } from "../types/project";
import { TranslationMatrix } from "../components/translation/TranslationTable";
import { ColumnProps } from "antd/es/table";

export const useTranslationData = (initialProjectId?: string) => {
  const [projects, setProjects] = useState<Project[]>([]);
  const [languages, setLanguages] = useState<Language[]>([]);
  const [columns, setColumns] = useState<ColumnProps[]>([]);
  const [selectedProject, setSelectedProject] = useState<number | null>(initialProjectId ? parseInt(initialProjectId) : null);
  const [loading, setLoading] = useState<boolean>(false);
  const [keyword, setKeyword] = useState<string>("");

  // Pagination state
  const [localPagination, setLocalPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  });

  const [paginatedMatrix, setPaginatedMatrix] = useState<TranslationMatrix[]>([]);

  // Selected items for batch operations
  const [selectedKeys, setSelectedKeys] = useState<string[]>([]);
  const [selectedTranslations, setSelectedTranslations] = useState<number[]>([]);
  const [batchDeleteLoading, setBatchDeleteLoading] = useState<boolean>(false);

  // 初始化
  useEffect(() => {
    fetchProjects();
    fetchLanguages();
  }, []);

  useEffect(() => {
    if (selectedProject) {
      fetchTranslations();
    }
  }, [selectedProject, keyword]);



  // 请求项目列表
  const fetchProjects = async () => {
    try {
      const data = await projectService.getProjects();
      setProjects(data);

      // 如果未选择项目，但有项目列表，则默认选择第一个项目
      if (!selectedProject && data.length > 0) {
        setSelectedProject(data[0].id);
      }
    } catch (error) {
      console.error("Failed to get project list:", error);
      message.error("Get project list failed");
    }
  };

  // 请求语言列表
  const fetchLanguages = async () => {
    try {
      const languages = await translationService.getLanguages();
      setLanguages(languages);
    } catch (error) {
      console.error("Failed to get language list:", error);
      message.error("Get language list failed");
    }
  };


  // 请求翻译列表
  const fetchTranslations = async () => {
    if (!selectedProject) return;

    try {
      setLoading(true);
      const result = await translationService.getTranslationMatrix(
        selectedProject,
        localPagination.current,
        localPagination.pageSize,
        keyword
      );

      setPaginatedMatrix(result.data);
      setLocalPagination({
        ...localPagination,
        total: result.total
      });
    } catch (error) {
      console.error("Failed to get translation matrix:", error);
      message.error("Get translation matrix failed");
    } finally {
      setLoading(false);
    }
  };

  // 表格分页
  const handleTableChange = (newPagination: TranslationTablePagination) => {
    const updatedPagination = {
      ...localPagination,
      current: newPagination.current,
      pageSize: newPagination.pageSize,
    };

    setLocalPagination(updatedPagination);
  };

  // 创建单条翻译
  const createTranslation = async (values: TranslationRequest) => {
    try {
      await translationService.createTranslation(values);
      message.success("Create translation successfully");
      fetchTranslations();
      return true;
    } catch (error) {
      console.error("Failed to create translation:", error);
      message.error("Create translation failed");
      return false;
    }
  };

  // 批量创建翻译
  const batchCreateTranslations = async (request: BatchTranslationRequest) => {
    try {
      await translationService.batchCreateTranslations(request);
      message.success("Batch create translations successfully");
      fetchTranslations();
      return true;
    } catch (error) {
      console.error("Failed to batch create translations:", error);
      message.error("Batch create translations failed");
      return false;
    }
  };

  // 更新翻译
  const updateTranslation = async (id: number, values: TranslationRequest) => {
    try {
      await translationService.updateTranslation(id, values);
      message.success("Update translation successfully");
      fetchTranslations();
      return true;
    } catch (error) {
      console.error("Failed to update translation:", error);
      message.error("Update translation failed");
      return false;
    }
  };

  // 删除翻译
  const deleteTranslation = async (id: number) => {
    try {
      await translationService.deleteTranslation(id);
      message.success("Delete translation successfully");
      fetchTranslations();
      return true;
    } catch (error) {
      console.error("Failed to delete translation:", error);
      message.error("Delete translation failed");
      return false;
    }
  };

  // 导出翻译
  const exportTranslations = async (format: string = "json") => {
    if (!selectedProject) {
      message.warning("Please select a project first");
      return null;
    }

    try {
      const data = await translationService.exportTranslations(
        selectedProject,
        format
      );
      message.success("Export translations successfully");
      return data;
    } catch (error) {
      console.error("Failed to export translations:", error);
      message.error("Export translations failed");
      return null;
    }
  };

  // 从JSON导入翻译
  const importTranslationsFromJson = async (data: Record<string, Record<string, string>>) => {
    if (!selectedProject) {
      message.warning("Please select a project first");
      return false;
    }

    try {
      await translationService.importTranslations(selectedProject, data);
      message.success("Import translations successfully");
      fetchTranslations();
      return true;
    } catch (error) {
      console.error("Failed to import translations:", error);
      message.error("Import translations failed");
      return false;
    }
  };

  // 批量删除翻译
  const batchDeleteTranslations = async () => {
    if (selectedTranslations.length === 0) {
      message.warning("Please select at least one translation record");
      return false;
    }

    try {
      setBatchDeleteLoading(true);
      await translationService.batchDeleteTranslations(selectedTranslations);

      message.success(`Successfully deleted ${selectedTranslations.length} translations`);
      setSelectedKeys([]);
      setSelectedTranslations([]);
      fetchTranslations();
      return true;
    } catch (error) {
      console.error("Failed to batch delete translations:", error);
      message.error("Batch delete translations failed");
      return false;
    } finally {
      setBatchDeleteLoading(false);
    }
  };

  // 处理行选择变化
  const handleRowSelectionChange = (
    selectedRowKeys: React.Key[],
    selectedRows: TranslationMatrix[]
  ) => {
    setSelectedKeys(selectedRowKeys as string[]);

    const selected: number[] = [];

    selectedRows.forEach((row) => {
      selected.push(...Object.values(row.languages).map(item => item.id))
    });

    setSelectedTranslations(selected);
  };

  useEffect(() => {
    if (selectedProject) {
      fetchTranslations();
    }
  }, [localPagination.current, localPagination.pageSize, selectedProject, keyword]);

  return {
    // State
    projects,
    languages,
    columns,
    selectedProject,
    loading,
    keyword,
    paginatedMatrix,
    localPagination,
    selectedKeys,
    selectedTranslations,
    batchDeleteLoading,

    // Actions
    setColumns,
    setSelectedProject,
    setKeyword,
    handleTableChange,
    fetchTranslations,
    createTranslation,
    batchCreateTranslations,
    updateTranslation,
    deleteTranslation,
    exportTranslations,
    importTranslationsFromJson,
    batchDeleteTranslations,
    handleRowSelectionChange,
  };
}; 