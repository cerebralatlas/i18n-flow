/**
 * Utility functions for working with localStorage
 */

// Keys for localStorage
export const STORAGE_KEYS = {
  LANGUAGE_COLUMNS: 'translation_language_columns',
};

/**
 * Save selected language columns to localStorage
 * @param columns Array of language codes to save
 */
export const saveSelectedColumns = (columns: string[]): void => {
  try {
    localStorage.setItem(STORAGE_KEYS.LANGUAGE_COLUMNS, JSON.stringify(columns));
  } catch (error) {
    console.error('Failed to save column preferences to localStorage:', error);
  }
};

/**
 * Load selected language columns from localStorage
 * @returns Array of language codes or null if not found
 */
export const loadSelectedColumns = (): string[] | null => {
  try {
    const storedColumns = localStorage.getItem(STORAGE_KEYS.LANGUAGE_COLUMNS);
    return storedColumns ? JSON.parse(storedColumns) : null;
  } catch (error) {
    console.error('Failed to load column preferences from localStorage:', error);
    return null;
  }
}; 