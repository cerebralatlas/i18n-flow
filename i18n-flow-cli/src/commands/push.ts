import { Command } from 'commander';
import fs from 'fs-extra';
import path from 'path';
import ora from 'ora';
import chalk from 'chalk';
import apiClient from '../core/api';
import { getConfig, validateConfig } from '../core/config';
import { logger } from '../utils/logger';
import { scanForTranslationKeys } from '../core/scanner';

export default function pushCommand(program: Command): void {
  program
    .command('push')
    .description('Push new translation keys to server')
    .option('-s, --scan', 'Scan source files for translation keys')
    .option('-p, --path <patterns>', 'Glob patterns for source files (comma-separated)')
    .option('-d, --dry-run', 'Show keys without pushing to server')
    .option('-n, --nested', 'Use nested directory structure (read from folders per locale)')
    .action(async (options) => {
      // 验证配置
      const { valid, missing } = validateConfig();
      if (!valid) {
        logger.error(`Missing configuration: ${missing.join(', ')}`);
        logger.info(`Please run ${chalk.cyan('i18n-flow init')} first`);
        return;
      }

      const config = getConfig();

      try {
        let keys: string[] = [];
        const defaults: Record<string, string> = {};

        // 扫描源码中的翻译键
        if (options.scan) {
          const spinner = ora('Scanning source files for translation keys...').start();

          // 使用自定义路径模式或配置中的默认值
          const patterns = options.path
            ? options.path.split(',')
            : config.sourcePatterns;

          // 扫描文件
          const result = await scanForTranslationKeys(patterns, config.extractorPattern);
          keys = result.keys;
          Object.assign(defaults, result.defaults);

          spinner.succeed(`Found ${keys.length} translation keys in source files`);
        } else if (options.nested) {
          // 从嵌套目录结构中读取翻译键
          const localesDir = path.resolve(process.cwd(), config.localesDir);
          const defaultLocaleDir = path.join(localesDir, config.defaultLocale);

          if (fs.existsSync(defaultLocaleDir) && fs.statSync(defaultLocaleDir).isDirectory()) {
            logger.info(`Reading keys from nested structure in ${config.defaultLocale} directory`);

            // 读取默认语言目录中的所有JSON文件
            const namespaceFiles = fs.readdirSync(defaultLocaleDir)
              .filter(file => file.endsWith('.json'));

            if (namespaceFiles.length === 0) {
              logger.warn(`No JSON files found in ${defaultLocaleDir}`);
              return;
            }

            // 处理每个命名空间文件
            for (const namespaceFile of namespaceFiles) {
              const namespace = namespaceFile.replace('.json', '');
              const filePath = path.join(defaultLocaleDir, namespaceFile);

              try {
                const namespaceTranslations = fs.readJSONSync(filePath);

                if (!namespaceTranslations || typeof namespaceTranslations !== 'object') {
                  logger.warn(`Invalid JSON in ${filePath}, skipping...`);
                  continue;
                }

                // 递归处理嵌套结构并添加命名空间前缀
                function processNestedObject(obj: any, currentPath: string) {
                  for (const [key, value] of Object.entries(obj)) {
                    const fullPath = currentPath ? `${currentPath}.${key}` : key;

                    if (value && typeof value === 'object' && !Array.isArray(value)) {
                      // 递归处理嵌套对象
                      processNestedObject(value, fullPath);
                    } else {
                      // 将完整路径的键添加到keys数组
                      const fullKey = `${namespace}.${fullPath}`;
                      keys.push(fullKey);
                      defaults[fullKey] = String(value);
                    }
                  }
                }

                processNestedObject(namespaceTranslations, '');
              } catch (error) {
                logger.warn(`Could not read ${filePath}: ${error instanceof Error ? error.message : 'Unknown error'}`);
              }
            }
            logger.info(`Loaded ${keys.length} keys from nested structure`);
          } else {
            logger.warn(`Default locale directory not found: ${defaultLocaleDir}`);
            logger.info('Make sure you have synced translations with --nested option first');
            return;
          }
        } else {
          // 从默认语言文件中获取键
          const defaultLocaleFile = path.join(
            process.cwd(),
            config.localesDir,
            `${config.defaultLocale}.json`
          );

          if (fs.existsSync(defaultLocaleFile)) {
            const defaultTranslations = fs.readJSONSync(defaultLocaleFile);
            keys = Object.keys(defaultTranslations);
            Object.assign(defaults, defaultTranslations);
            logger.info(`Loaded ${keys.length} keys from ${config.defaultLocale}.json`);
          } else {
            logger.warn(`Default locale file not found: ${defaultLocaleFile}`);
            logger.info('Use --scan option to scan source files for keys or --nested to read from nested structure');
            return;
          }
        }

        // 空检查
        if (keys.length === 0) {
          logger.warn('No translation keys found');
          return;
        }

        // 显示找到的键
        if (options.dryRun || keys.length <= 10) {
          logger.info('Translation keys:');
          keys.slice(0, 10).forEach(key => {
            const value = defaults[key] || '';
            logger.info(`  ${chalk.cyan(key)}: ${chalk.grey(value || '(no default)')}`);
          });

          if (keys.length > 10) {
            logger.info(`  ... and ${keys.length - 10} more`);
          }
        }

        // 仅展示模式
        if (options.dryRun) {
          logger.info(`Found ${keys.length} keys (dry run, not pushing to server)`);
          return;
        }

        // 推送到服务器
        const spinner = ora('Pushing translation keys to server...').start();

        const result = await apiClient.pushKeys({
          project_id: config.projectId,
          keys,
          defaults
        });

        spinner.succeed('Translation keys pushed to server');

        // 显示结果
        if (result.added?.length > 0) {
          logger.success(`Added ${result.added.length} new keys`);
        }

        if (result.existed?.length > 0) {
          result.existed.forEach(key => {
            logger.info(`Updated keys: ${chalk.green(key)}`);
          });
        }

        if (result.failed?.length > 0) {
          logger.warn(`Failed to add ${result.failed.length} keys`);
          logger.info('Failed keys:');
          result.failed.forEach(key => {
            logger.info(`  ${chalk.red(key)}`);
          });
        }
      } catch (error) {
        logger.error('Push failed');
        logger.error(error);
      }
    });
}