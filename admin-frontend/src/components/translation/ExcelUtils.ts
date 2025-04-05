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

      // Get the first sheet
      const firstSheetName = workbook.SheetNames[0];
      const worksheet = workbook.Sheets[firstSheetName];

      // Convert to JSON
      const jsonData = XLSX.utils.sheet_to_json(worksheet);

      if (jsonData.length === 0) {
        message.error("No data found in the Excel file");
        setLoading(false);
        callback(null);
        return;
      }

      // Generate preview table columns
      const sampleRow = jsonData[0] as Record<string, unknown>;
      const previewColumns = [
        {
          title: "Key Name",
          dataIndex: "key",
          key: "key",
          fixed: "left" as const,
          width: 200,
        },
      ];

      // Get all columns except the first one (assuming the first column is the key name) as possible language columns
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

      // Convert data to table available format
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
      console.error("Parse Excel file failed:", error);
      message.error("Parse Excel file failed, please ensure the file format is correct");
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
  // Convert Excel data to the format accepted by the backend
  // Format: { "zh-CN": { "key1": "value1", "key2": "value2" }, "en": { ... } }
  const importData: { [langCode: string]: { [key: string]: string } } = {};

  // Initialize language objects
  Object.values(selectedLanguages).forEach((langCode) => {
    importData[langCode] = {};
  });

  // Fill translation data
  excelData.forEach((row) => {
    const keyName = row["key"]; // Assuming the key column in Excel is 'key'
    if (!keyName) return; // Skip rows without key names

    // Add translations for each selected language column
    Object.entries(selectedLanguages).forEach(([excelCol, langCode]) => {
      if (row[excelCol] !== undefined && row[excelCol] !== null) {
        // Ensure all values are converted to strings
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