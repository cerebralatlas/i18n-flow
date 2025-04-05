import fs from 'fs-extra';
import path from 'path';
import { glob } from 'glob';

interface ScanResult {
  keys: string[];
  defaults: Record<string, string>;
  files: string[];
}

/**
 * 在源文件中扫描翻译键
 */
export async function scanForTranslationKeys(
  patterns: string[],
  extractorPattern: string
): Promise<ScanResult> {
  const result: ScanResult = {
    keys: [],
    defaults: {},
    files: []
  };

  // 编译正则表达式
  const regex = new RegExp(extractorPattern, 'g');

  // 查找匹配的文件
  for (const pattern of patterns) {
    const files = await glob.glob(pattern, { ignore: ['**/node_modules/**'] });
    result.files.push(...files);

    // 处理每个文件
    for (const file of files) {
      const content = fs.readFileSync(file, 'utf-8');

      // 匹配所有翻译键
      let match;
      while ((match = regex.exec(content)) !== null) {
        const key = match[1];
        if (key && !result.keys.includes(key)) {
          result.keys.push(key);

          // 尝试提取默认值 (如果有第二个参数并且是字符串字面量)
          try {
            // 简单的启发式方法，可能需要更复杂的AST解析
            const afterKey = content.substring(match.index + match[0].length);
            const defaultMatch = afterKey.match(/^\s*,\s*['"](.*?)['"][\),]/);
            if (defaultMatch) {
              result.defaults[key] = defaultMatch[1];
            }
          } catch (e) {
            // 默认值提取失败，忽略错误
          }
        }
      }
    }
  }

  return result;
}