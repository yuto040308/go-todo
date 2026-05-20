# go-todo

## 概要

Next.js + Go + PostgreSQL で構築した、シンプルな **ToDo 管理アプリ**です。

ToDo の CRUD（作成・一覧・更新・削除）をコア機能としつつ、**CI/CD パイプライン**や **AI コードレビュー**をはじめとするモダンな開発手法・最新技術を実践的に組み込むことを目的としたプロジェクトです。
「小さいアプリで、フルスタック開発のベストプラクティスをすべて体現する」ことをコンセプトに掲げています。

## デプロイ先 (Demo)

本プロジェクトは以下の環境にデプロイされており、実際に動作を確認できます。

- 🌐 **フロントエンド (Next.js)**: [https://go-todo-neon.vercel.app/](https://go-todo-neon.vercel.app/)
- ⚡ **バックエンド (Go/Gin)**: [https://go-todo-727829302986.europe-west1.run.app/api/hello](https://go-todo-727829302986.europe-west1.run.app/api/hello)（疎通確認用エンドポイント）

## 特徴

- シンプルな ToDo CRUD（作成・一覧・更新・削除）
- OpenAPI 仕様書をコードコメントから自動生成し、フロントエンドの型も自動生成
- **フロントエンド・バックエンド両方の静的解析を全面整備**（ESLint / Prettier / TypeScript型チェック / Knip / golangci-lint）
- **GitHub Actions による CI パイプライン**（フロント / バック 別ワークフロー、PR時に自動実行）
- **AI コードレビュー（PR-Agent / Gemini 2.5 Flash）** を PR フローに統合（PR説明文の自動生成・レビュー・改善提案）
- Vercel / Cloud Run / Supabase を活用したサーバーレス・マネージドな本番環境

---

## 技術構成

### フロントエンド

| カテゴリ | 技術 | バージョン | 用途 |
| --- | --- | --- | --- |
| フレームワーク | Next.js (App Router) | 16.2.4 | Web フレームワーク |
| ランタイム | React | 19.2.4 | UI ライブラリ |
| 言語 | TypeScript | 5.x | 型安全な JavaScript |
| スタイリング | Tailwind CSS | v4 | ユーティリティファーストの CSS |
| UI コンポーネント | shadcn/ui (radix-ui) | - | Tailwind ベースのコンポーネント集 |
| データフェッチ | TanStack Query | v5.100.3 | サーバー状態管理・キャッシュ |
| 型自動生成 | openapi-typescript | v7.13.0 | OpenAPI 仕様書から TypeScript 型を生成 |

> 最新の細かいバージョンは [frontend/package.json](frontend/package.json) を参照。

---

### バックエンド

| カテゴリ | 技術 | バージョン | 用途 |
| --- | --- | --- | --- |
| 言語 | Go | 1.25.0 | バックエンド言語 |
| フレームワーク | Gin | v1.12.0 | HTTP Web フレームワーク |
| API ドキュメント | swaggo/swag | v1.16.6 | コードコメントから OpenAPI 仕様書（yaml）を自動生成 |

> swaggo/swag は v2.0.0-rc5（RC版）もありますが、安定版の v1.16.6 を採用しています。
> 最新の細かいバージョンは [backend/go.mod](backend/go.mod) を参照。

---

### データベース

| カテゴリ | 技術 | バージョン | 用途 |
| --- | --- | --- | --- |
| DB | PostgreSQL | 17.x | リレーショナルデータベース |

> Supabase が対応しているバージョンに合わせて 17.x を採用しています。

---

### インフラ

| サービス | 用途 | デプロイ先URL |
| --- | --- | --- |
| Vercel | Next.js フロントエンドのホスティング・デプロイ | [go-todo-neon.vercel.app](https://go-todo-neon.vercel.app/) |
| Cloud Run | Go / Gin API サーバーのホスティング（Google Cloud） | [go-todo-727829302986.europe-west1.run.app](https://go-todo-727829302986.europe-west1.run.app/) |
| Supabase | PostgreSQL のホスティング・認証・ストレージ | （※接続設定のみ） |

---

### 品質保証ツール

コード品質を担保するため、以下のツールを導入しています。

| ツール | 対象 | 役割 |
| --- | --- | --- |
| ESLint | フロント | JavaScript / TypeScript の構文・バグ検出 |
| Prettier | フロント | コードフォーマッタ（見た目の統一） |
| eslint-config-prettier | フロント | ESLint と Prettier の競合を解消 |
| TypeScript (`tsc --noEmit`) | フロント | 型整合性のチェック |
| Knip | フロント | 未使用ファイル・export・依存パッケージの検出 |
| golangci-lint | バック | Go の静的解析メタリンター（errcheck / staticcheck / govet 等を内包） |
| Reviewdog | CI共通 | 静的解析結果を PR に「Commit suggestion」付きで自動提案 |
| PR-Agent | CI共通 | AI コードレビュー（Gemini 2.5 Flash 利用） |

---

## ローカル開発

### 起動

```bash
docker compose up -d --build
# または
make up
```

### アクセス

ブラウザからは **必ず Nginx 経由（ポート80）でアクセス**してください。

| URL | 用途 |
| --- | --- |
| `http://localhost` | ✅ Nginx経由（推奨）。`/api/*` がバックエンドに振り分けられる |
| `http://localhost:3000` | ❌ Next.js直アクセス。`/api/*` が解決できず 404 になる |

`http://localhost:3000` でアクセスすると、フロントから `/api/hello` を叩いてもNext.js自身が応答してしまい、Goバックエンドに到達しません。Nginxを経由することでフロント・APIが同一オリジンに統一され、CORSも回避できます。

---

## コード品質チェック

ローカルでもCIと同じ静的解析を再現できるよう、**Makefile に全コマンドを集約**しています。
詳細は後述の「Makefile チートシート」を参照。

### フロントエンド静的解析

フロントエンドは4種類の静的解析を組み合わせています。

| ツール | チェックコマンド | 自動修正コマンド | 設定ファイル |
| --- | --- | --- | --- |
| **ESLint** | `make lint-frontend` | `make lint-fix-frontend` | [frontend/eslint.config.mjs](frontend/eslint.config.mjs) |
| **Prettier** | `make format-check-frontend` | `make format-frontend` | [frontend/.prettierrc](frontend/.prettierrc) / [frontend/.prettierignore](frontend/.prettierignore) |
| **TypeScript 型チェック** | `make typecheck-frontend` | （手動修正） | [frontend/tsconfig.json](frontend/tsconfig.json) |
| **Knip（未使用検出）** | `make unused-check-frontend` | （手動削除） | [frontend/knip.json](frontend/knip.json) |

#### ESLint と Prettier の役割分担

- **ESLint**：未使用変数・React フック違反など「**バグや品質**」を検出
- **Prettier**：インデント・クォート・改行など「**見た目**」を整える
- 両者の競合を防ぐため [`eslint-config-prettier`](https://github.com/prettier/eslint-config-prettier) を ESLint 設定の最後に適用しています（フォーマット系ルールを ESLint 側で無効化）

#### Knip について

Knip は「**使われていないファイル・export・依存パッケージ**」を検出します。
shadcn/ui のスキャフォールドファイルや `openapi-typescript` 等は誤検知となるため、[frontend/knip.json](frontend/knip.json) で除外設定済みです。

### バックエンド静的解析

バックエンドコンテナには **golangci-lint** が同梱されています。
Goコードの静的解析・コーディング規約チェックが実行できます。

| 用途 | コマンド |
| --- | --- |
| 静的解析を実行 | `make lint` |
| 静的解析 + 自動修正 | `make lint-fix` |

> golangci-lint は `errcheck` `staticcheck` `govet` など複数のlinterを束ねたメタリンターで、デフォルト設定でも実用十分な指摘が得られます。

---

## CI/CD パイプライン

PR を作成すると以下のワークフローが自動実行されます。**CIが赤いままだとマージ不可**な運用にすることで、品質を担保しています。

### ワークフロー一覧

| ワークフロー | トリガー | 内容 |
| --- | --- | --- |
| [frontend_ci_lint.yml](.github/workflows/frontend_ci_lint.yml) | PR (`frontend/**` 変更時) | ESLint / TypeScript 型チェック / Prettier チェック / Knip |
| [golangci_lint.yml](.github/workflows/golangci_lint.yml) | PR | golangci-lint + Reviewdog による自動提案 |
| [go_test.yml](.github/workflows/go_test.yml) | PR (`backend/**` 変更時) | `go test ./...` でユニットテスト実行 |
| [pr_agent.yml](.github/workflows/pr_agent.yml) | PR opened/reopened/ready_for_review/synchronize | AI コードレビュー（PR-Agent） |

### Reviewdog による静的解析の自動提案

`golangci_lint.yml` では Reviewdog を使用し、PR 上に直接結果を表示します。

| 指摘の種類 | PR上での見え方 | 対応方法 |
| --- | --- | --- |
| **自動修正できるもの**（gofumpt の整形など） | 「Commit suggestion」**ボタン付きの提案コメント** | ボタンをクリックすると、修正コミットがそのまま PR ブランチに push される |
| **自動修正できないもの**（errcheck の `_ = err` など） | 行単位のインラインコメント（ボタンなし） | 自分でコードを直して push し直す |

### AI コードレビュー（PR-Agent）

[pr_agent.yml](.github/workflows/pr_agent.yml) では、Google の **Gemini 2.5 Flash** を使用した AI レビューを自動実行します。
全アクションが**日本語で出力**されるよう設定済み。

| アクション | 内容 |
| --- | --- |
| **describe** | PR の説明文（タイトル・要約・変更点リスト）を自動生成 |
| **review** | シニアエンジニア視点でコードレビューを実施し、設計・保守性の改善提案 |
| **improve** | 具体的なコード改善案を提示 |

---

## Makefile チートシート

プロジェクトルートに `Makefile` を用意しているので、長いコマンドを覚えなくても `make xxx` で簡単に操作できます。

### 全体操作

| コマンド | 内容 |
| --- | --- |
| `make up` | 全コンテナを起動 |
| `make down` | 全コンテナを停止 |
| `make rebuild-backend` | バックエンド（Go）をビルドし直して起動 |
| `make rebuild-frontend` | フロントエンド（Next.js）をビルドし直して起動 |
| `make rebuild-nginx` | Nginx をビルドし直して起動 |

### バックエンド

| コマンド | 内容 |
| --- | --- |
| `make lint` | Go の静的解析（golangci-lint）を実行 |
| `make lint-fix` | Go の静的解析 + 自動修正 |
| `make test` | Go のテスト（`go test ./...`）を実行 |

### フロントエンド

| コマンド | 内容 |
| --- | --- |
| `make lint-frontend` | ESLint を実行 |
| `make lint-fix-frontend` | ESLint を実行 + 自動修正 |
| `make typecheck-frontend` | TypeScript 型チェック（`tsc --noEmit`） |
| `make format-frontend` | Prettier を一括適用（**ファイル書き換えあり**） |
| `make format-check-frontend` | Prettier ルール通りに整形されているかチェック（書き換えなし、CI 想定） |
| `make unused-check-frontend` | Knip で未使用コード・依存パッケージを検出 |
| `make frontend-install` | フロントエンドの依存パッケージをインストール（引数: `pkg="package-name"`） |
| `make reset-frontend` | フロントエンドを匿名ボリュームごと作り直す（`node_modules` が壊れた時用） |

### 使用例

```bash
# コンテナをまとめて起動
make up

# バックエンドのコードを変えた後、ビルドし直して起動
make rebuild-backend

# バックエンドの静的解析を実行
make lint

# フロントエンドの全静的解析を順番に実行（CI と同じチェック）
make lint-frontend
make typecheck-frontend
make format-check-frontend
make unused-check-frontend

# フロントエンドに Prettier を一括適用
make format-frontend

# フロントエンドに新しいパッケージをインストール
make frontend-install pkg="zustand"
```

> `./...` は「カレントディレクトリ以下のすべてのGoパッケージを対象」というGo標準の表記です。
