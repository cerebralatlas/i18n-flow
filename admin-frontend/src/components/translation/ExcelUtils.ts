import * as XLSX from "xlsx";
import { Language } from "../../types/translation";
import { message } from "antd";

export interface ExcelData {
  jsonData: any[];
  columns: any[];
  previewData: any[];
}

/**
 * Parse Excel file and prepare data for preview
 */
export const parseExcelFile = (
  file: File,
  languages: Language[],
  callback: (data: ExcelData | null) => void,
  setLoading: (loading: boolean) => void
) => {
  setLoading(true);

  const reader = new FileReader();
  reader.onload = (e: ProgressEvent<FileReader>) => {
    try {
      if (!e.target || !e.target.result) {
        message.error("读取文件失败");
        setLoading(false);
        callback(null);
        return;
      }

      const data = new Uint8Array(e.target.result as ArrayBuffer);
      const workbook = XLSX.read(data, { type: "array" });

      // 获取第一个工作表
      const firstSheetName = workbook.SheetNames[0];
      const worksheet = workbook.Sheets[firstSheetName];

      // 转换为JSON
      const jsonData = XLSX.utils.sheet_to_json(worksheet);

      if (jsonData.length === 0) {
        message.error("Excel文件中没有找到数据");
        setLoading(false);
        callback(null);
        return;
      }

      // 生成预览表格列
      const sampleRow = jsonData[0] as Record<string, unknown>;
      const previewColumns = [
        {
          title: "键名",
          dataIndex: "key",
          key: "key",
          fixed: "left" as const,
          width: 200,
        },
      ];

      // 获取除了第一列(假设第一列是键名)外的所有列作为可能的语言列
      const possibleLanguageColumns = Object.keys(sampleRow).filter(
        (key) => key !== "key"
      );

      possibleLanguageColumns.forEach((colName) => {
        previewColumns.push({
          title: colName,
          dataIndex: colName,
          key: colName,
          width: 200,
          fixed: "left" as const,
        });
      });

      // 转换数据为表格可用格式
      const previewData = jsonData.map((row, index) => ({
        key: index.toString(),
        ...row as Record<string, unknown>,
      }));

      callback({
        jsonData,
        columns: previewColumns,
        previewData,
      });
    } catch (error) {
      console.error("解析Excel文件失败:", error);
      message.error("解析Excel文件失败，请确保文件格式正确");
      callback(null);
    } finally {
      setLoading(false);
    }
  };

  reader.readAsArrayBuffer(file);
};

/**
 * Format Excel data for import
 */
export const formatExcelDataForImport = (
  excelData: any[],
  selectedLanguages: { [key: string]: string }
) => {
  // 将Excel数据转换为后端接受的格式
  // 格式: { "zh-CN": { "key1": "value1", "key2": "value2" }, "en": { ... } }
  const importData: { [langCode: string]: { [key: string]: string } } = {};

  // 初始化语言对象
  Object.values(selectedLanguages).forEach((langCode) => {
    importData[langCode] = {};
  });

  // 填充翻译数据
  excelData.forEach((row) => {
    const keyName = row["key"]; // 假设Excel中的键名列是 'key'
    if (!keyName) return; // 跳过没有键名的行

    // 为每个选定的语言列添加翻译
    Object.entries(selectedLanguages).forEach(([excelCol, langCode]) => {
      if (row[excelCol] !== undefined && row[excelCol] !== null) {
        // 确保所有值都转换为字符串
        importData[langCode][keyName] = String(row[excelCol]);
      }
    });
  });

  return importData;
};

/**
 * Auto-map Excel columns to languages
 */
export const autoMapLanguageColumns = (
  columns: any[],
  languages: Language[]
) => {
  const initialLanguageMapping: { [key: string]: string } = {};

  // Only process columns after the first one (which is the key name)
  columns.slice(1).forEach((column) => {
    const colName = column.title;

    // Try to match with system languages
    const matchedLanguage = languages.find(
      (lang) =>
        colName.toLowerCase().includes(lang.code.toLowerCase()) ||
        colName.toLowerCase().includes(lang.name.toLowerCase())
    );

    if (matchedLanguage) {
      initialLanguageMapping[colName] = matchedLanguage.code;
    }
  });

  return initialLanguageMapping;
}; 