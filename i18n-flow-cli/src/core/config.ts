import fs from 'fs-extra';
import path from 'path';
import os from 'os';
import Conf from 'conf';

// 配置文件名
const CONFIG_FILENAME = '.i18nflowrc.json';

// 全局配置存储
const globalConfig = new Conf({ projectName: 'i18n-flow' });

// 配置类型定义
export interface I18nFlowConfig {
  serverUrl: string;
  apiKey: string;
  projectId: string;
  localesDir: string;
  defaultLocale: string;
  sourcePatterns: string[];
  extractorPattern: string;
  languageMapping?: Record<string, string>; // 语言代码映射：文件名 -> 数据库语言代码
}

// 默认配置
const defaultConfig: I18nFlowConfig = {
  serverUrl: 'http://localhost:8080',
  apiKey: '',
  projectId: '',
  localesDir: './src/locales',
  defaultLocale: 'en',
  sourcePatterns: [
    'src/**/*.{js,jsx,ts,tsx}',
    '!src/**/*.{spec,test}.{js,jsx,ts,tsx}',
    '!**/node_modules/**'
  ],
  extractorPattern: "(?:t|i18n\\.t)\\(['\"]([\\w\\.\\-]+)['\"]"
};

/**
 * 获取项目根目录
 */
export function findProjectRoot(startDir = process.cwd()): string | null {
  let currentDir = startDir;

  // 向上查找直到找到package.json或根目录
  while (currentDir !== path.parse(currentDir).root) {
    if (fs.existsSync(path.join(currentDir, 'package.json'))) {
      return currentDir;
    }
    currentDir = path.dirname(currentDir);
  }

  return null;
}

/**
 * 获取配置文件路径
 */
export function getConfigPath(): string {
  const projectRoot = findProjectRoot();
  if (!projectRoot) {
    throw new Error('Could not find project root directory');
  }
  return path.join(projectRoot, CONFIG_FILENAME);
}

/**
 * 检查配置文件是否存在
 */
export function configExists(): boolean {
  return fs.existsSync(getConfigPath());
}

/**
 * 获取配置
 */
export function getConfig(): I18nFlowConfig {
  // 首先尝试从项目配置获取
  const configPath = getConfigPath();
  if (fs.existsSync(configPath)) {
    const projectConfig = fs.readJSONSync(configPath);
    return { ...defaultConfig, ...projectConfig };
  }

  // 如果没有项目配置，返回默认配置
  return defaultConfig;
}

/**
 * 保存配置
 */
export function saveConfig(config: Partial<I18nFlowConfig>): void {
  const configPath = getConfigPath();
  const currentConfig = configExists() ? getConfig() : defaultConfig;
  const newConfig = { ...currentConfig, ...config };

  fs.writeJSONSync(configPath, newConfig, { spaces: 2 });
}

/**
 * 获取API密钥
 */
export function getApiKey(): string {
  const config = getConfig();
  return config.apiKey;
}

/**
 * 保存API密钥
 */
export function saveApiKey(apiKey: string): void {
  saveConfig({ apiKey });
}

/**
 * 验证配置是否完整
 */
export function validateConfig(): { valid: boolean; missing: string[] } {
  const config = getConfig();
  const missing: string[] = [];

  if (!config.apiKey) missing.push('apiKey');
  if (!config.projectId) missing.push('projectId');
  if (!config.serverUrl) missing.push('serverUrl');

  return {
    valid: missing.length === 0,
    missing
  };
}