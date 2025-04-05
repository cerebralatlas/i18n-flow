export interface TranslationMap {
  [key: string]: string;
}

export interface LocaleTranslations {
  [locale: string]: TranslationMap;
}

export interface ProjectConfig {
  id: string;
  name: string;
  description: string;
}

export interface ConflictEntry {
  key: string;
  localValue: string;
  remoteValue: string;
  resolution?: string;
}

export interface ConflictResult {
  conflicts: ConflictEntry[];
  resolved: boolean;
}

export interface ScanOptions {
  patterns: string[];
  extractorPattern: string;
  ignorePatterns?: string[];
}