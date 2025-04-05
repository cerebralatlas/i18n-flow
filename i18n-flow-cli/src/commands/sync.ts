import { Command } from 'commander';
import fs from 'fs-extra';
import path from 'path';
import ora from 'ora';
import chalk from 'chalk';
import apiClient from '../core/api';
import { getConfig, validateConfig } from '../core/config';
import { logger } from '../utils/logger';

export default function syncCommand(program: Command): void {
  program
    .command('sync')
    .description('Sync translations from server')
    .option('-l, --locale <locales>', 'Comma-separated list of locales to sync')
    .option('-f, --force', 'Force overwrite local translations')
    .option('-n, --nested', 'Use nested directory structure (one folder per locale, split by first key part)')
    .action(async (options) => {
      // 验证配置
      const { valid, missing } = validateConfig();
      if (!valid) {
        logger.error(`Missing configuration: ${missing.join(', ')}`);
        logger.info(`Please run ${chalk.cyan('i18n-flow init')} first`);
        return;
      }

      const config = getConfig();
      const spinner = ora('Syncing translations').start();

      try {
        // 解析要同步的语言
        const locales = options.locale ? options.locale.split(',') : [];

        // 从服务器获取翻译
        const translations = await apiClient.getTranslations(config.projectId);

        // 准备本地文件路径
        const localesDir = path.resolve(process.cwd(), config.localesDir);
        fs.ensureDirSync(localesDir);

        // 获取翻译中的语言列表
        const availableLocales = new Set<string>();
        Object.values(translations).forEach(translationObj => {
          Object.keys(translationObj as Record<string, any>).forEach(locale => {
            availableLocales.add(locale);
          });
        });

        // 过滤要处理的语言
        const localesToProcess = locales.length > 0
          ? Array.from(availableLocales).filter(locale => locales.includes(locale))
          : Array.from(availableLocales);

        if (localesToProcess.length === 0) {
          spinner.warn('No translations found for the specified locales');
          return;
        }

        // 按语言处理翻译
        for (const locale of localesToProcess) {
          // 构建该语言的翻译对象
          const localeTranslations: Record<string, string> = {};

          // 填充翻译数据
          Object.entries(translations).forEach(([key, translationObj]) => {
            const localeValue = (translationObj as Record<string, any>)[locale];
            if (localeValue) {
              localeTranslations[key] = localeValue;
            }
          });

          if (options.nested) {
            // 使用嵌套目录结构
            const localeDir = path.join(localesDir, locale);
            fs.ensureDirSync(localeDir);

            // 根据key的第一部分分组
            const groupedTranslations: Record<string, Record<string, string>> = {};

            Object.entries(localeTranslations).forEach(([key, value]) => {
              const [namespace, ...rest] = key.split('.');
              if (!namespace) return;

              if (!groupedTranslations[namespace]) {
                groupedTranslations[namespace] = {};
              }

              // 剩余部分作为新的key
              const newKey = rest.join('.');
              if (newKey) {
                groupedTranslations[namespace][newKey] = value;
              }
            });

            // 保存每个分组到独立文件
            for (const [namespace, namespaceTranslations] of Object.entries(groupedTranslations)) {
              const namespaceFilePath = path.join(localeDir, `${namespace}.json`);

              // 处理已存在的文件
              if (fs.existsSync(namespaceFilePath) && !options.force) {
                // 合并现有翻译
                const existingTranslations = fs.readJSONSync(namespaceFilePath, { throws: false }) || {};
                const mergedTranslations = { ...existingTranslations, ...namespaceTranslations };
                fs.writeJSONSync(namespaceFilePath, mergedTranslations, { spaces: 2 });
              } else {
                // 创建新文件
                fs.writeJSONSync(namespaceFilePath, namespaceTranslations, { spaces: 2 });
              }
            }
          } else {
            // 使用传统的单文件结构
            const localeFilePath = path.join(localesDir, `${locale}.json`);

            // 处理已存在的文件
            if (fs.existsSync(localeFilePath) && !options.force) {
              // 合并现有翻译
              const existingTranslations = fs.readJSONSync(localeFilePath, { throws: false }) || {};
              const mergedTranslations = { ...existingTranslations, ...localeTranslations };
              fs.writeJSONSync(localeFilePath, mergedTranslations, { spaces: 2 });
            } else {
              // 创建新文件
              fs.writeJSONSync(localeFilePath, localeTranslations, { spaces: 2 });
            }
          }
        }

        spinner.succeed(`Synced translations for ${localesToProcess.length} locales`);
        logger.info(`Translations saved to: ${chalk.cyan(localesDir)}`);

        if (options.nested) {
          logger.info(`Used nested structure: ${chalk.cyan('Each locale has its own folder with namespace files')}`);
        }

        // 统计信息
        const keyCount = Object.keys(translations).length;
        logger.info(`Total keys: ${chalk.cyan(keyCount.toString())}`);
        logger.info(`Locales: ${chalk.cyan(localesToProcess.join(', '))}`);
      } catch (error) {
        spinner.fail('Sync failed');
        logger.error(error);
      }
    });
}