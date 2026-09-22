import js from "@eslint/js";
import tseslint from "typescript-eslint";
import globals from "globals";
import pluginReact from "eslint-plugin-react";
import reactHooks from "eslint-plugin-react-hooks";
import reactRefresh from "eslint-plugin-react-refresh";
import simpleImportSort from "eslint-plugin-simple-import-sort";
import prettier from "eslint-config-prettier";
import { defineConfig } from "eslint/config";

const sourceFiles = ["src/**/*.{js,jsx,ts,tsx}"];
const reactFiles = ["src/**/*.{jsx,tsx}"];
const configFiles = ["vite.config.ts", "eslint.config.js"];

export default defineConfig([
  {
    ignores: ["dist/**", "coverage/**", "src-tauri/**", "src/routeTree.gen.ts"],
  },

  // React アプリ本体
  {
    files: sourceFiles,
    plugins: {
      js,
      "simple-import-sort": simpleImportSort,
    },
    extends: ["js/recommended", tseslint.configs.recommended],
    languageOptions: {
      globals: globals.browser,
    },
    rules: {
      "simple-import-sort/imports": "error",
      "simple-import-sort/exports": "error",
    },
  },

  // Vite / ESLint 設定
  {
    files: configFiles,
    plugins: {
      js,
    },
    extends: ["js/recommended", tseslint.configs.recommended],
    languageOptions: {
      globals: globals.node,
    },
  },

  // React の基本ルール
  {
    ...pluginReact.configs.flat.recommended,
    files: reactFiles,
    settings: {
      react: {
        version: "19.1",
      },
    },
  },

  // React 19 の JSX Transform 対応
  {
    ...pluginReact.configs.flat["jsx-runtime"],
    files: reactFiles,
  },

  // Hooks のルール
  {
    ...reactHooks.configs.flat.recommended,
    files: reactFiles,
  },

  // Vite の Fast Refresh 対応
  {
    ...reactRefresh.configs.vite,
    files: reactFiles,
  },

  {
    files: ["src/routes/**/*.{jsx,tsx}"],
    rules: {
      "react-refresh/only-export-components": "off",
    },
  },

  // Prettier と競合する ESLint ルールを無効化
  prettier,
]);
