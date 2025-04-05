import chalk from 'chalk';

export const logger = {
  info: (message: string): void => {
    console.log(`${chalk.blue('ℹ')} ${message}`);
  },

  success: (message: string): void => {
    console.log(`${chalk.green('✓')} ${message}`);
  },

  warn: (message: string): void => {
    console.log(`${chalk.yellow('⚠')} ${message}`);
  },

  error: (error: any): void => {
    if (error instanceof Error) {
      console.error(`${chalk.red('✗')} ${error.message}`);
      if (process.env.DEBUG) {
        console.error(error.stack);
      }
    } else {
      console.error(`${chalk.red('✗')} ${error}`);
    }
  },

  debug: (message: string): void => {
    if (process.env.DEBUG) {
      console.log(`${chalk.gray('🔧')} ${message}`);
    }
  }
};