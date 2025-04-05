import { useState, useEffect } from "react";
import { message } from "antd";
import { translationService } from "../services/translationService";
import { projectService } from "../services/projectService";
import {
  TranslationResponse,
  Language,
  BatchTranslationRequest,
} from "../types/translation";
import { Project } from "../types/project";
import { TranslationMatrix } from "../components/translation/TranslationTable";

export const useTranslationData = (initialProjectId?: string) => {
  // State variables
  const [translations, setTranslations] = useState<TranslationResponse[]>([]);
  const [translationMatrix, setTranslationMatrix] = useState<TranslationMatrix[]>([]);
  const [projects, setProjects] = useState<Project[]>([]);
  const [languages, setLanguages] = useState<Language[]>([]);
  const [columns, setColumns] = useState<any[]>([]);
  const [selectedProject, setSelectedProject] = useState<number | null>(initialProjectId ? parseInt(initialProjectId) : null);
  const [loading, setLoading] = useState<boolean>(false);
  const [keyword, setKeyword] = useState<string>("");

  // Pagination state
  const [localPagination, setLocalPagination] = useState({
    current: 1,
    pageSize: 10,
    total: 0,
  });

  const [allTranslations, setAllTranslations] = useState<TranslationResponse[]>([]);
  const [paginatedMatrix, setPaginatedMatrix] = useState<TranslationMatrix[]>([]);

  // Selected items for batch operations
  const [selectedKeys, setSelectedKeys] = useState<string[]>([]);
  const [selectedTranslations, setSelectedTranslations] = useState<TranslationResponse[]>([]);
  const [batchDeleteLoading, setBatchDeleteLoading] = useState<boolean>(false);

  // Initialize
  useEffect(() => {
    fetchProjects();
    fetchLanguages();
  }, []);

  useEffect(() => {
    if (selectedProject) {
      fetchTranslations();
    }
  }, [selectedProject, keyword]);

  // When local pagination parameters change, recalculate data to display
  useEffect(() => {
    if (translationMatrix.length > 0) {
      const { current, pageSize } = localPagination;
      const startIndex = (current - 1) * pageSize;
      const endIndex = startIndex + pageSize;
      setPaginatedMatrix(translationMatrix.slice(startIndex, endIndex));
    }
  }, [localPagination.current, localPagination.pageSize, translationMatrix]);

  // Get project list
  const fetchProjects = async () => {
    try {
      const { data } = await projectService.getProjects();
      setProjects(data);

      // If no project is selected but there is a project list, select the first one by default
      if (!selectedProject && data.length > 0) {
        setSelectedProject(data[0].id);
      }
    } catch (error) {
      console.error("Failed to get project list:", error);
      message.error("获取项目列表失败");
    }
  };

  // Get language list
  const fetchLanguages = async () => {
    try {
      const languages = await translationService.getLanguages();
      setLanguages(languages);
    } catch (error) {
      console.error("Failed to get language list:", error);
      message.error("获取语言列表失败");
    }
  };



  // Get translation list
  const fetchTranslations = async () => {
    if (!selectedProject) return;

    try {
      setLoading(true);
      const response = await translationService.getTranslationMatrix(
        selectedProject,
        localPagination.current,
        localPagination.pageSize,
        keyword
      );

      setPaginatedMatrix(response.data.data);
      setLocalPagination({
        ...localPagination,
        total: response.data.total
      });
    } catch (error) {
      console.error("Failed to get translation matrix:", error);
      message.error("获取翻译矩阵失败");
    } finally {
      setLoading(false);
    }
  };

  // Handle pagination change
  const handleTableChange = (newPagination: any) => {
    const updatedPagination = {
      ...localPagination,
      current: newPagination.current,
      pageSize: newPagination.pageSize,
    };

    setLocalPagination(updatedPagination);
  };

  // Create a single translation
  const createTranslation = async (values: any) => {
    try {
      await translationService.createTranslation(values);
      message.success("创建翻译成功");
      fetchTranslations();
      return true;
    } catch (error) {
      console.error("Failed to create translation:", error);
      message.error("创建翻译失败");
      return false;
    }
  };

  // Batch create translations
  const batchCreateTranslations = async (request: BatchTranslationRequest) => {
    try {
      await translationService.batchCreateTranslations(request);
      message.success("批量创建翻译成功");
      fetchTranslations();
      return true;
    } catch (error) {
      console.error("Failed to batch create translations:", error);
      message.error("批量创建翻译失败");
      return false;
    }
  };

  // Update translation
  const updateTranslation = async (id: number, values: any) => {
    try {
      await translationService.updateTranslation(id, values);
      message.success("更新翻译成功");
      fetchTranslations();
      return true;
    } catch (error) {
      console.error("Failed to update translation:", error);
      message.error("更新翻译失败");
      return false;
    }
  };

  // Delete translation
  const deleteTranslation = async (id: number) => {
    try {
      await translationService.deleteTranslation(id);
      message.success("删除翻译成功");
      fetchTranslations();
      return true;
    } catch (error) {
      console.error("Failed to delete translation:", error);
      message.error("删除翻译失败");
      return false;
    }
  };

  // Export translations
  const exportTranslations = async (format: string = "json") => {
    if (!selectedProject) {
      message.warning("请先选择项目");
      return null;
    }

    try {
      const data = await translationService.exportTranslations(
        selectedProject,
        format
      );
      message.success("导出翻译成功");
      return data;
    } catch (error) {
      console.error("Failed to export translations:", error);
      message.error("导出翻译失败");
      return null;
    }
  };

  // Import translations from JSON
  const importTranslationsFromJson = async (data: any) => {
    if (!selectedProject) {
      message.warning("请先选择项目");
      return false;
    }

    try {
      await translationService.importTranslations(selectedProject, data);
      message.success("导入翻译成功");
      fetchTranslations();
      return true;
    } catch (error) {
      console.error("Failed to import translations:", error);
      message.error("导入翻译失败");
      return false;
    }
  };

  // Batch delete translations
  const batchDeleteTranslations = async () => {
    if (selectedTranslations.length === 0) {
      message.warning("请至少选择一条翻译记录");
      return false;
    }

    try {
      setBatchDeleteLoading(true);
      // Extract IDs of selected translations
      const ids = selectedTranslations.map((item) => item.id);

      // Call batch delete API
      await translationService.batchDeleteTranslations(ids);

      message.success(`成功删除 ${ids.length} 条翻译`);
      setSelectedKeys([]);
      setSelectedTranslations([]);
      fetchTranslations();
      return true;
    } catch (error) {
      console.error("Failed to batch delete translations:", error);
      message.error("批量删除翻译失败");
      return false;
    } finally {
      setBatchDeleteLoading(false);
    }
  };

  // Handle row selection change
  const handleRowSelectionChange = (
    selectedRowKeys: React.Key[],
    selectedRows: TranslationMatrix[]
  ) => {
    setSelectedKeys(selectedRowKeys as string[]);

    // Find all corresponding translation records based on selected rows
    const selected: TranslationResponse[] = [];

    selectedRows.forEach((row) => {
      // For each key name, collect all language translation records from allTranslations
      const keyTranslations = allTranslations.filter(
        (t) => t.key_name === row.key_name
      );
      selected.push(...keyTranslations);
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
    translations,
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