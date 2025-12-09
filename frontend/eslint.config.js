import js from "@eslint/js";
import globals from "globals";
import pluginVue from "eslint-plugin-vue";
import css from "@eslint/css";
import security from "eslint-plugin-security";
import noSecrets from "eslint-plugin-no-secrets";
import sonarjs from 'eslint-plugin-sonarjs';
import { defineConfig } from "eslint/config";

export default defineConfig([
  { files: ["**/*.{js,mjs,cjs,vue}"], plugins: { js }, extends: ["js/recommended"], languageOptions: { globals: globals.browser } },
  // Vue files
  ...pluginVue.configs["flat/essential"].map(cfg => ({
    ...cfg,
    files: ["**/*.vue"], 
  })),
  {
    files: ["**/*.{js,vue}"],
    plugins: {
      "no-secrets": noSecrets,
    },
    rules: {
      "no-secrets/no-secrets": "error",
    },
  },
  {
    files: ["**/*.js", "**/*.vue"],
    ...sonarjs.configs.recommended,
  },
  {
    files: ["**/*.js", "**/*.vue"],
    ...security.configs.recommended,
  },
  { files: ["**/*.css"], plugins: { css }, language: "css/css", extends: ["css/recommended"] },
  {
		// Note: there should be no other properties in this object
		ignores: ["dist/*", "node_modules/*"],
	},
]);
