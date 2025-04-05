import fs from 'fs-extra';
import path from 'path';
import { findProjectRoot } from '../core/config';

/**
 * 深度合并两个对象
 */
export function deepMerge<T extends Record<string, any>>(target: T, source: Record<string, any>): T {
  const output = { ...target };

  if (isObject(target) && isObject(source)) {
    Object.keys(source).forEach(key => {
      if (isObject(source[key])) {
        if (!(key in target)) {
          Object.assign(output, { [key]: source[key] });
        } else {
          // 使用类型断言确保TypeScript理解这是安全的操作
          output[key as keyof T] = deepMerge(target[key] as Record<string, any>, source[key]) as T[keyof T];
        }
      } else {
        Object.assign(output, { [key]: source[key] });
      }
    });
  }

  return output;
}

/**
 * 检查值是否为对象
 */
function isObject(item: any): item is Record<string, any> {
  return (item && typeof item === 'object' && !Array.isArray(item));
}

/**
 * 获取包的版本号
 */
export function getPackageVersion(): string {
  try {
    const projectRoot = findProjectRoot(__dirname);
    if (!projectRoot) return '1.0.0';

    const packageJsonPath = path.join(projectRoot, 'package.json');
    const packageJson = fs.readJSONSync(packageJsonPath) as { version?: string };
    return packageJson.version || '1.0.0';
  } catch (error) {
    return '1.0.0';
  }
}

/**
 * 递归创建嵌套的翻译对象
 * 
 * 例如: { 'user.login.title': 'Login' } => { user: { login: { title: 'Login' } } }
 */
export function createNestedTranslations(translations: Record<string, string>): Record<string, any> {
  const result: Record<string, any> = {};

  Object.entries(translations).forEach(([key, value]) => {
    const parts = key.split('.');
    let current = result;

    for (let i = 0; i < parts.length - 1; i++) {
      const part = parts[i];
      if (!current[part]) {
        current[part] = {};
      }
      current = current[part];
    }

    current[parts[parts.length - 1]] = value;
  });

  return result;
}

/**
 * 将嵌套的翻译对象扁平化
 * 
 * 例如: { user: { login: { title: 'Login' } } } => { 'user.login.title': 'Login' }
 */
export function flattenTranslations(obj: Record<string, any>, prefix = ''): Record<string, string> {
  const result: Record<string, string> = {};

  for (const key in obj) {
    const prefixedKey = prefix ? `${prefix}.${key}` : key;

    if (typeof obj[key] === 'object' && obj[key] !== null) {
      Object.assign(result, flattenTranslations(obj[key], prefixedKey));
    } else {
      result[prefixedKey] = obj[key];
    }
  }

  return result;
}