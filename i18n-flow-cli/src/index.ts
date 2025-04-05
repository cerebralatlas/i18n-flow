#!/usr/bin/env node
import { Command } from 'commander';
import chalk from 'chalk';
import initCommand from './commands/init';
import syncCommand from './commands/sync';
import pushCommand from './commands/push';
import statusCommand from './commands/status';
import { getPackageVersion } from './utils/file';

// 创建程序实例
const program = new Command();

// 设置版本号和描述
const version = getPackageVersion();
program
  .name('i18n-flow')
  .description('i18n-flow CLI for translation management')
  .version(version);

// 注册命令
initCommand(program);
syncCommand(program);
pushCommand(program);
statusCommand(program);

// 帮助信息
program.on('--help', () => {
  console.log('');
  console.log(`Run ${chalk.cyan('i18n-flow <command> --help')} for detailed usage of given command.`);
  console.log('');
  console.log('Example:');
  console.log(`  ${chalk.cyan('i18n-flow init')}          Initialize i18n-flow in your project`);
  console.log(`  ${chalk.cyan('i18n-flow sync')}          Sync translations from server`);
  console.log(`  ${chalk.cyan('i18n-flow push --scan')}   Scan and push new keys to server`);
});

// 解析命令行参数
program.parse(process.argv);

// 如果没有提供命令，则显示帮助信息
if (!process.argv.slice(2).length) {
  program.outputHelp();
}