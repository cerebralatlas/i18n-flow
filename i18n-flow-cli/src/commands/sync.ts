import { Command } from 'commander';
import fs from 'fs-extra';
import path from 'path';
import ora from 'ora';
import chalk from 'chalk';
import apiClient from '../core/api';
import { getConfig, validateConfig } from '../core/config';
import { logger } from '../utils/logger';

/**
 * 将扁平化的键值对转换为嵌套对象结构
 */
function convertFlatToNested(flatObject: Record<string, string>): Record<string, any> {
  const result: Record<string, any> = {};
  
  for (const [key, value] of Object.entries(flatObject)) {
    const keys = key.split('.');
    let current = result;
    
    // 遍历除最后一个键之外的所有键，创建嵌套对象
    for (let i = 0; i < keys.length - 1; i++) {
      const keyPart = keys[i];
      if (!current[keyPart] || typeof current[keyPart] !== 'object') {
        current[keyPart] = {};
      }
      current = current[keyPart];
    }
    
    // 设置最终值
    const lastKey = keys[keys.length - 1];
    current[lastKey] = value;
  }
  
  return result;
}

/**
 * 深度合并两个对象
 */
function mergeDeep(target: Record<string, any>, source: Record<string, any>): Record<string, any> {
  const result = { ...target };
  
  for (const [key, value] of Object.entries(source)) {
    if (value && typeof value === 'object' && !Array.isArray(value)) {
      if (result[key] && typeof result[key] === 'object' && !Array.isArray(result[key])) {
        result[key] = mergeDeep(result[key], value);
      } else {
        result[key] = value;
      }
    } else {
      result[key] = value;
    }
  }
  
  return result;
}

export default function syncCommand(program: Command): void {
  program
    .command('sync')
    .description('Sync translations from server')
    .option('-l, --locale <locales>', 'Comma-separated list of locales to sync')
    .option('-f, --force', 'Force overwrite local translations')
    .option('-n, --nested', 'Use nested directory structure (one folder per locale, split by first key part)')
    .action(async (options) => {
      // validate the config
      const { valid, missing } = validateConfig();
      if (!valid) {
        logger.error(`Missing configuration: ${missing.join(', ')}`);
        logger.info(`Please run ${chalk.cyan('i18n-flow init')} first`);
        return;
      }

      const config = getConfig();
      const spinner = ora('Syncing translations').start();

      try {
        // parse the locales to sync
        const locales = options.locale ? options.locale.split(',') : [];

        // get the translations from the server
        const translations = await apiClient.getTranslations(config.projectId);

        // prepare the local file path
        const localesDir = path.resolve(process.cwd(), config.localesDir);
        fs.ensureDirSync(localesDir);

        // get the list of languages from the translations
        const availableLocales = new Set<string>();
        
        // Check if we have any translation keys
        const translationKeys = Object.keys(translations);
        if (translationKeys.length > 0) {
          // Look at the first translation to get available languages
          const firstTranslation = translations[translationKeys[0]];
          if (typeof firstTranslation === 'object' && firstTranslation !== null) {
            Object.keys(firstTranslation).forEach(locale => {
              availableLocales.add(locale);
            });
          }
        }

        // filter the locales to process
        const localesToProcess = locales.length > 0
          ? Array.from(availableLocales).filter(locale => locales.includes(locale))
          : Array.from(availableLocales);

        if (localesToProcess.length === 0) {
          spinner.warn('No translations found for the specified locales');
          return;
        }

        // process translations by language
        for (const locale of localesToProcess) {
          // build the translation object for the language
          const localeTranslations: Record<string, string> = {};

          // fill the translation data
          Object.entries(translations).forEach(([key, translationObj]) => {
            const localeValue = (translationObj as Record<string, any>)[locale];
            if (localeValue) {
              localeTranslations[key] = localeValue;
            }
          });

          if (options.nested) {
            // use nested directory structure - create one folder per locale with namespace files
            const localeDir = path.join(localesDir, locale);
            fs.ensureDirSync(localeDir);

            // group by the first part of the key (namespace)
            const groupedTranslations: Record<string, Record<string, any>> = {};

            Object.entries(localeTranslations).forEach(([key, value]) => {
              const keyParts = key.split('.');
              if (keyParts.length < 2) {
                // 如果键名没有点分隔，跳过或放入一个默认命名空间
                logger.warn(`Skipping key without namespace: ${key}`);
                return;
              }

              const namespace = keyParts[0];
              if (!groupedTranslations[namespace]) {
                groupedTranslations[namespace] = {};
              }

              // 构建嵌套对象结构
              let current = groupedTranslations[namespace];
              for (let i = 1; i < keyParts.length - 1; i++) {
                const part = keyParts[i];
                if (!current[part]) {
                  current[part] = {};
                }
                current = current[part];
              }

              // 设置最终值
              const lastPart = keyParts[keyParts.length - 1];
              current[lastPart] = value;
            });

            // save each namespace to a separate file
            for (const [namespace, namespaceTranslations] of Object.entries(groupedTranslations)) {
              const namespaceFilePath = path.join(localeDir, `${namespace}.json`);

              // handle existing files
              if (fs.existsSync(namespaceFilePath) && !options.force) {
                // merge the existing translations
                const existingTranslations = fs.readJSONSync(namespaceFilePath, { throws: false }) || {};
                const mergedTranslations = { ...existingTranslations, ...namespaceTranslations };
                fs.writeJSONSync(namespaceFilePath, mergedTranslations, { spaces: 2 });
              } else {
                // create new file
                fs.writeJSONSync(namespaceFilePath, namespaceTranslations, { spaces: 2 });
              }
            }
          } else {
            // use the traditional single file structure with nested objects
            const localeFilePath = path.join(localesDir, `${locale}.json`);
            
            // 将扁平化的键转换为嵌套对象结构
            const nestedTranslations = convertFlatToNested(localeTranslations);

            // handle existing files
            if (fs.existsSync(localeFilePath) && !options.force) {
              // merge the existing translations
              const existingTranslations = fs.readJSONSync(localeFilePath, { throws: false }) || {};
              const mergedTranslations = mergeDeep(existingTranslations, nestedTranslations);
              fs.writeJSONSync(localeFilePath, mergedTranslations, { spaces: 2 });
            } else {
              // create new file
              fs.writeJSONSync(localeFilePath, nestedTranslations, { spaces: 2 });
            }
          }
        }

        spinner.succeed(`Synced translations for ${localesToProcess.length} locales`);
        logger.info(`Translations saved to: ${chalk.cyan(localesDir)}`);

        if (options.nested) {
          logger.info(`Used nested structure: ${chalk.cyan('Each locale has its own folder with namespace files')}`);
        }

        // statistics
        const keyCount = Object.keys(translations).length;
        logger.info(`Total keys: ${chalk.cyan(keyCount.toString())}`);
        logger.info(`Locales: ${chalk.cyan(localesToProcess.join(', '))}`);
      } catch (error) {
        spinner.fail('Sync failed');
        logger.error(error);
      }
    });
}