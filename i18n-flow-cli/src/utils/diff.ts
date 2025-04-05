/**
 * 比较两个翻译集合，返回差异
 */
export function compareTranslations(
  source: Record<string, string>,
  target: Record<string, string>
): {
  added: string[];
  modified: string[];
  removed: string[];
  unchanged: string[];
} {
  const sourceKeys = Object.keys(source);
  const targetKeys = Object.keys(target);

  const added = targetKeys.filter(key => !sourceKeys.includes(key));
  const removed = sourceKeys.filter(key => !targetKeys.includes(key));

  const common = sourceKeys.filter(key => targetKeys.includes(key));
  const modified = common.filter(key => source[key] !== target[key]);
  const unchanged = common.filter(key => source[key] === target[key]);

  return {
    added,
    modified,
    removed,
    unchanged
  };
}

/**
 * 生成两个翻译集合的合并报告
 */
export function generateMergeReport(
  local: Record<string, string>,
  remote: Record<string, string>
): {
  local: Record<string, string>;
  remote: Record<string, string>;
  conflicts: Record<string, { local: string; remote: string }>;
} {
  const diff = compareTranslations(local, remote);
  const conflicts: Record<string, { local: string; remote: string }> = {};

  // 找出修改冲突
  diff.modified.forEach(key => {
    conflicts[key] = {
      local: local[key],
      remote: remote[key]
    };
  });

  return {
    local: diff.added.reduce((acc: Record<string, string>, key: string) => {
      acc[key] = remote[key];
      return acc;
    }, {}),
    remote: diff.removed.reduce((acc: Record<string, string>, key: string) => {
      acc[key] = local[key];
      return acc;
    }, {}),
    conflicts
  };
}