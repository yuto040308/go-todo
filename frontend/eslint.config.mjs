import { defineConfig, globalIgnores } from 'eslint/config';
import nextVitals from 'eslint-config-next/core-web-vitals';
import nextTs from 'eslint-config-next/typescript';
import prettier from 'eslint-config-prettier';

const eslintConfig = defineConfig([
  ...nextVitals,
  ...nextTs,
  // Override default ignores of eslint-config-next.
  globalIgnores([
    // Default ignores of eslint-config-next:
    '.next/**',
    'out/**',
    'build/**',
    'next-env.d.ts',
  ]),
  // ESLintとPrettierはフォーマット系ルールで衝突する可能性がある
  // 配列の最後に置くことで、ESLint側のフォーマットルールを全OFFにして、
  // 見た目はPrettierに一任する
  prettier,
]);

export default eslintConfig;
