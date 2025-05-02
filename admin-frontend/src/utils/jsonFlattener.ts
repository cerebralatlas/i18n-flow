// 定义嵌套的翻译数据类型
interface NestedTranslations {
  [key: string]: string | NestedTranslations;
}

// 定义扁平化后的翻译数据类型
interface FlattenedTranslations {
  [key: string]: string;
}

// 定义多语言翻译数据类型
interface LanguageTranslations {
  [language: string]: NestedTranslations;
}

/**
 * 将嵌套的对象扁平化为点号分隔的键值对
 * @param obj 嵌套的对象
 * @param prefix 键前缀
 * @returns 扁平化后的对象
 */
export function flattenJSON(
  obj: NestedTranslations,
  prefix: string = ''
): FlattenedTranslations {
  const flattened: FlattenedTranslations = {};

  for (const key in obj) {
    const value = obj[key];
    if (typeof value === 'object' && value !== null && !Array.isArray(value)) {
      const nested = flattenJSON(
        value as NestedTranslations,
        prefix ? `${prefix}.${key}` : key
      );
      Object.assign(flattened, nested);
    } else if (typeof value === 'string') {
      flattened[prefix ? `${prefix}.${key}` : key] = value;
    }
  }

  return flattened;
}

/**
 * 处理多语言翻译JSON数据
 * @param data 多语言翻译数据
 * @returns 处理后的扁平化多语言翻译数据
 */
export function processTranslationJSON(
  data: LanguageTranslations
): Record<string, FlattenedTranslations> {
  const result: Record<string, FlattenedTranslations> = {};

  for (const lang in data) {
    result[lang] = flattenJSON(data[lang]);
  }

  return result;
} 