import { Command } from 'commander';
import inquirer from 'inquirer';
import chalk from 'chalk';
import fs from 'fs-extra';
import path from 'path';
import ora from 'ora';
import apiClient from '../core/api';
import { saveConfig, configExists, getConfig } from '../core/config';
import { logger } from '../utils/logger';

export default function initCommand(program: Command): void {
  program
    .command('init')
    .description('Initialize i18n-flow in your project')
    .option('-y, --yes', 'Skip confirmation prompts and use default values')
    .action(async (options) => {
      logger.info('Initializing i18n-flow CLI...');

      if (configExists() && !options.yes) {
        const { confirm } = await inquirer.prompt([
          {
            type: 'confirm',
            name: 'confirm',
            message: 'i18n-flow config already exists. Do you want to overwrite it?',
            default: false
          }
        ]);

        if (!confirm) {
          logger.info('Initialization cancelled.');
          return;
        }
      }

      // Get inputs from user
      const answers = await inquirer.prompt([
        {
          type: 'input',
          name: 'serverUrl',
          message: 'Enter i18n-flow server URL:',
          default: getConfig().serverUrl || 'http://localhost:8080'
        },
        {
          type: 'input',
          name: 'apiKey',
          message: 'Enter API key:',
          validate: (input) => input.length > 0 ? true : 'API key is required'
        },
        {
          type: 'input',
          name: 'projectId',
          message: 'Enter project ID:',
          validate: (input) => input.length > 0 ? true : 'Project ID is required'
        },
        {
          type: 'input',
          name: 'localesDir',
          message: 'Enter locales directory:',
          default: getConfig().localesDir || './src/locales'
        },
        {
          type: 'input',
          name: 'defaultLocale',
          message: 'Enter default locale:',
          default: getConfig().defaultLocale || 'en'
        }
      ]);

      // Test connection
      const spinner = ora('Testing server connection...').start();

      try {
        // Temporarily save API key for testing
        saveConfig({ apiKey: answers.apiKey, serverUrl: answers.serverUrl });

        const connected = await apiClient.testConnection();

        if (!connected) {
          spinner.fail('Could not connect to the server. Please check your API key and server URL.');
          return;
        }

        spinner.succeed('Connected to the server successfully!');

        // Save config
        saveConfig({
          serverUrl: answers.serverUrl,
          apiKey: answers.apiKey,
          projectId: answers.projectId,
          localesDir: answers.localesDir,
          defaultLocale: answers.defaultLocale
        });

        // Ensure locales directory exists
        const localesPath = path.resolve(process.cwd(), answers.localesDir);
        fs.ensureDirSync(localesPath);

        logger.success('i18n-flow initialized successfully!');
        logger.info(`Locales directory: ${chalk.cyan(localesPath)}`);
        logger.info(`Next step: Run ${chalk.cyan('i18n-flow sync')} to download translations`);
      } catch (error) {
        spinner.fail('Initialization failed');
        logger.error(error);
      }
    });
}