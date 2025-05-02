import i18n from "i18next";
import { initReactI18next } from "react-i18next";
// 导入翻译文件
import enTranslations from "../messages/en.json";
import zhTranslations from "../messages/zh.json";

// the translations
// (tip move them in a JSON file and import them,
// or even better, manage them separated from your code: https://react.i18next.com/guides/multiple-translation-files)
const resources = {
  en: {
    translation: enTranslations
  },
  zh: {
    translation: zhTranslations
  }
};

i18n
  .use(initReactI18next) // passes i18n down to react-i18next
  .init({
    resources,
    lng: "en", // 默认语言
    fallbackLng: "en", // 当某个翻译键缺失时，使用的备选语言
    interpolation: {
      escapeValue: false
    }
  });

export default i18n;