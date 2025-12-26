// Translation data models
export interface Translation {
  id: number
  project_id: number
  language_id: number
  key_name: string
  value: string
  context?: string
  status: 'active' | 'deprecated'
  created_at: string
  updated_at: string
}

export interface Language {
  id: number
  code: string
  name: string
  is_default: boolean
  status: 'active' | 'inactive'
  created_at: string
  updated_at: string
}

// Translation matrix for table display
export interface TranslationMatrixRow {
  key_name: string
  context?: string
  translations: Record<string, TranslationCell>
}

export interface TranslationCell {
  id?: number
  language_id: number
  value: string
  status?: 'active' | 'deprecated'
}

export interface TranslationMatrix {
  languages: Language[]
  rows: TranslationMatrixRow[]
  total_count: number
  page: number
  page_size: number
  total_pages: number
}

// API request models
export interface CreateTranslationRequest {
  project_id: number
  language_id: number
  key_name: string
  value: string
  context?: string
}

export interface BatchTranslationRequest {
  project_id: number
  key_name: string
  context?: string
  translations: Record<string, string> // language_code -> value
}

export interface CreateLanguageRequest {
  code: string
  name: string
  is_default?: boolean
}

export interface ImportTranslationsData {
  [languageCode: string]: {
    [key: string]: string
  }
}
