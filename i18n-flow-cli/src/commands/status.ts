import { Command } from 'commander';
import fs from 'fs-extra';
import path from 'path';
import chalk from 'chalk';
import ora from 'ora';
import apiClient from '../core/api';
import { getConfig, validateConfig } from '../core/config';
import { logger } from '../utils/logger';

export default function statusCommand(program: Command): void {
  program
    .command('status')
    .description('Check translation status')
    .action(async () => {
      // 验证配置
      const { valid, missing } = validateConfig();
      if (!valid) {
        logger.error(`Missing configuration: ${missing.join(', ')}`);
        logger.info(`Please run ${chalk.cyan('i18n-flow init')} first`);
        return;
      }

      const config = getConfig();
      const spinner = ora('Checking translation status...').start();

      try {
        // 从服务器获取翻译
        const serverTranslations = await apiClient.getTranslations(config.projectId);

        // 获取本地翻译文件
        const localesDir = path.resolve(process.cwd(), config.localesDir);
        if (!fs.existsSync(localesDir)) {
          spinner.info('Locales directory does not exist');
          return;
        }

        // 读取本地翻译文件
        const localeFiles = fs.readdirSync(localesDir)
          .filter(file => file.endsWith('.json'))
          .map(file => file.replace('.json', ''));

        // 服务器上的键总数
        const serverKeys = Object.keys(serverTranslations);
        spinner.succeed(`Found ${serverKeys.length} keys on server and ${localeFiles.length} locale files locally`);

        // 检查每个语言的状态
        for (const locale of localeFiles) {
          const localeFilePath = path.join(localesDir, `${locale}.json`);
          const localTranslations = fs.readJSONSync(localeFilePath);
          const localKeys = Object.keys(localTranslations);

          // 对比服务器和本地的键
          const missingLocally = serverKeys.filter(key => !localKeys.includes(key));
          const extraLocally = localKeys.filter(key => !serverKeys.includes(key));

          logger.info(`${chalk.cyan(locale)}: ${localKeys.length} keys`);

          if (missingLocally.length > 0) {
            logger.warn(`  - Missing ${missingLocally.length} keys from server`);
            if (missingLocally.length <= 5) {
              missingLocally.forEach(key => {
                logger.info(`    ${chalk.yellow(key)}`);
              });
            }
          }

          if (extraLocally.length > 0) {
            logger.warn(`  - Extra ${extraLocally.length} keys not on server`);
            if (extraLocally.length <= 5) {
              extraLocally.forEach(key => {
                logger.info(`    ${chalk.gray(key)}`);
              });
            }
          }

          if (missingLocally.length === 0 && extraLocally.length === 0) {
            logger.success(`  - In sync with server`);
          }
        }

        // 建议
        if (localeFiles.length === 0) {
          logger.info(`Tip: Run ${chalk.cyan('i18n-flow sync')} to download translations from server`);
        } else if (serverKeys.length === 0) {
          logger.info(`Tip: Run ${chalk.cyan('i18n-flow push --scan')} to push keys to server`);
        }
      } catch (error) {
        spinner.fail('Status check failed');
        logger.error(error);
      }
    });
}