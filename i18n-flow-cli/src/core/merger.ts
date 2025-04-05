import fs from 'fs-extra';
import path from 'path';
import { compareTranslations, generateMergeReport } from '../utils/diff';
import { TranslationMap } from '../types';
import { logger } from '../utils/logger';

/**
 * 合并本地和远程翻译
 */
export async function mergeTranslations(
  localPath: string,
  remoteTranslations: TranslationMap,
  options: {
    force?: boolean;
    backup?: boolean;
    conflictFile?: string;
  } = {}
): Promise<{
  merged: boolean;
  conflicts: number;
  added: number;
  removed: number;
  unchanged: number;
}> {
  // 读取本地翻译
  let localTranslations: TranslationMap = {};
  if (fs.existsSync(localPath)) {
    try {
      localTranslations = await fs.readJSON(localPath);
    } catch (error) {
      logger.warn(`Could not parse ${localPath}, treating as empty`);
    }
  }

  // 如果强制模式，直接用远程覆盖
  if (options.force) {
    await fs.writeJSON(localPath, remoteTranslations, { spaces: 2 });
    return {
      merged: true,
      conflicts: 0,
      added: Object.keys(remoteTranslations).length,
      removed: 0,
      unchanged: 0
    };
  }

  // 创建备份
  if (options.backup && fs.existsSync(localPath)) {
    const backupPath = `${localPath}.backup-${Date.now()}`;
    await fs.copy(localPath, backupPath);
    logger.info(`Created backup at ${backupPath}`);
  }

  // 比较差异
  const diff = compareTranslations(localTranslations, remoteTranslations);

  // 生成合并报告
  const report = generateMergeReport(localTranslations, remoteTranslations);

  // 处理冲突
  if (Object.keys(report.conflicts).length > 0 && options.conflictFile) {
    const conflictPath = options.conflictFile;
    await fs.writeJSON(
      conflictPath,
      {
        timestamp: new Date().toISOString(),
        conflicts: Object.entries(report.conflicts).map(([key, values]) => ({
          key,
          localValue: values.local,
          remoteValue: values.remote
        }))
      },
      { spaces: 2 }
    );
    logger.warn(`Conflicts found and saved to ${conflictPath}`);
  }

  // 合并翻译
  const mergedTranslations = {
    ...localTranslations,
    ...remoteTranslations
  };

  // 写入合并后的文件
  await fs.writeJSON(localPath, mergedTranslations, { spaces: 2 });

  return {
    merged: true,
    conflicts: Object.keys(report.conflicts).length,
    added: diff.added.length,
    removed: 0, // 保持本地的移除键
    unchanged: diff.unchanged.length
  };
}