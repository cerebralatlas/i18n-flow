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
        Object.values(translations).forEach(translationObj => {
          Object.keys(translationObj as Record<string, any>).forEach(locale => {
            availableLocales.add(locale);
          });
        });

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
            // use nested directory structure
            const localeDir = path.join(localesDir, locale);
            fs.ensureDirSync(localeDir);

            // group by the first part of the key
            const groupedTranslations: Record<string, Record<string, any>> = {};

            Object.entries(localeTranslations).forEach(([key, value]) => {
              const [namespace, ...rest] = key.split('.');
              if (!namespace) return;

              if (!groupedTranslations[namespace]) {
                groupedTranslations[namespace] = {};
              }

              // rest part become new key
              const newKey = rest.join('.');
              if (newKey) {
                // For keys like "common.signinSuccess" (simple case)
                if (!newKey.includes('.')) {
                  groupedTranslations[namespace][newKey] = value;
                } else {
                  // For keys like "banner.info.title" (nested case)
                  const parts = newKey.split('.');
                  let current = groupedTranslations[namespace];

                  // Create nested objects for all parts except the last one
                  for (let i = 0; i < parts.length - 1; i++) {
                    const part = parts[i];
                    if (!current[part]) {
                      current[part] = {};
                    }
                    current = current[part];
                  }

                  // Set the value at the deepest level
                  const lastPart = parts[parts.length - 1];
                  current[lastPart] = value;
                }
              }
            });

            // save each group to a separate file
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
            // use the traditional single file structure
            const localeFilePath = path.join(localesDir, `${locale}.json`);

            // handle existing files
            if (fs.existsSync(localeFilePath) && !options.force) {
              // merge the existing translations
              const existingTranslations = fs.readJSONSync(localeFilePath, { throws: false }) || {};
              const mergedTranslations = { ...existingTranslations, ...localeTranslations };
              fs.writeJSONSync(localeFilePath, mergedTranslations, { spaces: 2 });
            } else {
              // create new file
              fs.writeJSONSync(localeFilePath, localeTranslations, { spaces: 2 });
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