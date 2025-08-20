/**
 * 导出工具函数
 * 支持多种格式的文件导出：JSON、Excel、CSV等
 */

// 下载文件的通用函数
export const downloadFile = (content: string | Blob, filename: string, contentType: string) => {
  const blob = content instanceof Blob ? content : new Blob([content], { type: contentType });
  const url = URL.createObjectURL(blob);
  const link = document.createElement('a');
  link.href = url;
  link.download = filename;
  document.body.appendChild(link);
  link.click();
  document.body.removeChild(link);
  URL.revokeObjectURL(url);
};

// 导出JSON格式
export const exportToJSON = (data: Record<string, Record<string, string>>, projectId: number) => {
  const jsonStr = JSON.stringify(data, null, 2);
  const filename = `translations_${projectId}_${new Date().toISOString().split('T')[0]}.json`;
  downloadFile(jsonStr, filename, 'application/json');
};

// 导出CSV格式
export const exportToCSV = (data: Record<string, Record<string, string>>, projectId: number) => {
  // 获取所有语言代码
  const allLanguages = new Set<string>();
  Object.values(data).forEach(translations => {
    Object.keys(translations).forEach(lang => allLanguages.add(lang));
  });
  const languages = Array.from(allLanguages).sort();

  // 创建CSV头部
  const headers = ['Key', ...languages];
  const csvRows = [headers.join(',')];

  // 添加数据行
  Object.entries(data).forEach(([key, translations]) => {
    const row = [
      `"${key.replace(/"/g, '""')}"`, // 转义双引号
      ...languages.map(lang => {
        const value = translations[lang] || '';
        return `"${value.replace(/"/g, '""')}"`;
      })
    ];
    csvRows.push(row.join(','));
  });

  const csvContent = csvRows.join('\n');
  const filename = `translations_${projectId}_${new Date().toISOString().split('T')[0]}.csv`;
  downloadFile(csvContent, filename, 'text/csv');
};

// 导出Excel格式（使用SheetJS）
export const exportToExcel = async (data: Record<string, Record<string, string>>, projectId: number) => {
  try {
    // 动态导入SheetJS
    const XLSX = await import('xlsx');
    
    // 获取所有语言代码
    const allLanguages = new Set<string>();
    Object.values(data).forEach(translations => {
      Object.keys(translations).forEach(lang => allLanguages.add(lang));
    });
    const languages = Array.from(allLanguages).sort();

    // 创建工作表数据
    const worksheetData: (string | number)[][] = [];
    
    // 添加标题行
    worksheetData.push(['Key', ...languages]);
    
    // 添加数据行
    Object.entries(data).forEach(([key, translations]) => {
      const row: (string | number)[] = [key, ...languages.map(lang => translations[lang] || '')];
      worksheetData.push(row);
    });

    // 创建工作簿和工作表
    const workbook = XLSX.utils.book_new();
    const worksheet = XLSX.utils.aoa_to_sheet(worksheetData);
    
    // 设置列宽
    const columnWidths = [
      { wch: 30 }, // Key列宽度
      ...languages.map(() => ({ wch: 20 })) // 语言列宽度
    ];
    worksheet['!cols'] = columnWidths;

    // 添加工作表到工作簿
    XLSX.utils.book_append_sheet(workbook, worksheet, 'Translations');

    // 导出文件
    const filename = `translations_${projectId}_${new Date().toISOString().split('T')[0]}.xlsx`;
    XLSX.writeFile(workbook, filename);
    
  } catch (error) {
    console.error('Excel export failed:', error);
    throw new Error('Excel导出失败，请确保已安装相关依赖');
  }
};

// 统一的导出函数
export const exportTranslations = (
  data: Record<string, Record<string, string>>, 
  format: 'json' | 'csv' | 'excel', 
  projectId: number
) => {
  switch (format) {
    case 'json':
      exportToJSON(data, projectId);
      break;
    case 'csv':
      exportToCSV(data, projectId);
      break;
    case 'excel':
      exportToExcel(data, projectId);
      break;
    default:
      throw new Error(`不支持的导出格式: ${format}`);
  }
};