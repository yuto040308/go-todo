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
- **スキーマ駆動開発（spec-first）**：`api/openapi.yaml` を Single Source of Truth とし、`make gen-api` で Go の型・Gin サーバインタフェース・TypeScript の型を一括自動生成
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
| 型自動生成 | openapi-typescript | v7.13.0 | `api/openapi.yaml` から TypeScript 型を生成（→ `frontend/types/api.ts`） |

> 最新の細かいバージョンは [frontend/package.json](frontend/package.json) を参照。

---

### バックエンド

| カテゴリ | 技術 | バージョン | 用途 |
| --- | --- | --- | --- |
| 言語 | Go | 1.25.0 | バックエンド言語 |
| フレームワーク | Gin | v1.12.0 | HTTP Web フレームワーク |
| ORM | GORM | v1.31.x | Go から PostgreSQL を操作 |
| 認証 | golang-jwt/jwt | v5.x | `Authorization: Bearer <JWT>` 形式の認証ミドルウェア |
| 型自動生成 | oapi-codegen | v2.7.x | `api/openapi.yaml` から Go の型 + Gin サーバインタフェースを生成（→ `backend/gen/api/`） |
| Swagger UI | swaggest/swgui | v1.8.x | ローカル開発時のみ `/swagger/` で OpenAPI 仕様を可視化 |

> 最新の細かいバージョンは [backend/go.mod](backend/go.mod) を参照。

---

### データベース

| カテゴリ | 技術 | バージョン | 用途 |
| --- | --- | --- | --- |
| DB | PostgreSQL | 18.x | リレーショナルデータベース |
| マイグレーション | [golang-migrate](https://github.com/golang-migrate/migrate) | v4.19.x | バージョン管理された SQL マイグレーション（Docker イメージ経由で実行） |

> ローカル開発では Docker コンテナで PostgreSQL を起動します。本番 (Supabase) の対応バージョンに合わせる必要があれば、Supabase 接続チケットで調整します。

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

## OpenAPI / スキーマ駆動開発

API の型は `api/openapi.yaml` を **Single Source of Truth** とし、そこから Go と TypeScript の型を一括自動生成します。

### ワークフロー

```
api/openapi.yaml                       ← ここを編集する（人間が書く）
    │
    │ make gen-api
    ▼
┌──────────────────────────┐   ┌──────────────────────────┐
│ backend/gen/api/         │   │ frontend/types/api.ts    │
│  └ api.gen.go            │   │  (paths / components)    │
│    (型 + ServerInterface)│   │                          │
└──────────────────────────┘   └──────────────────────────┘
    │                              │
    ▼                              ▼
  Gin handler が実装           axios / TanStack Query
  ServerInterface を満たす       のレスポンス型に使う
```

### 関連コマンド

| コマンド | 内容 |
| --- | --- |
| `make gen-api` | Go と TS の型をまとめて再生成 |
| `make gen-api-backend` | Go の型のみ再生成（`backend/gen/api/api.gen.go`） |
| `make gen-api-frontend` | TS の型のみ再生成（`frontend/types/api.ts`） |

### Swagger UI（ローカル限定）

backend 起動時、`APP_ENV` が `production` 以外なら以下で OpenAPI 仕様をブラウザで確認できます。

| URL | 内容 |
| --- | --- |
| `http://localhost:8080/swagger/` | Swagger UI（インタラクティブな仕様ビューア） |
| `http://localhost:8080/openapi.yaml` | yaml 自体の配信（UI が裏で取得） |

> 本番環境（Cloud Run）では `APP_ENV=production` を設定することで自動的に無効化されます。

### 生成物の取り扱い

`backend/gen/api/api.gen.go` と `frontend/types/api.ts` は **commit します**（`.gitignore` しない）。理由：

- レビュー時に「yaml の変更で型がどう変わったか」が PR の diff で確認できる
- CI で「コミット済みの生成物が yaml と一致しているか」を `git diff --exit-code` で検出できる（仕様逸脱の検出）
- 新規 clone 時に毎回 `make gen-api` を走らせる必要が無くなる
- エディタの go-to-definition が型ファイルに飛べる

これらのファイルは **直接編集しない**（先頭に `Code generated ... DO NOT EDIT.` の注意書きあり）。型を変えたいときは必ず `api/openapi.yaml` を編集して `make gen-api` を実行してください。

---

## API エンドポイント

すべて `/api` 配下。認証は `Authorization: Bearer <JWT>` 方式（詳細は [api/openapi.yaml](api/openapi.yaml) が SoT）。

### 認証 (auth)

| メソッド | パス | 認証 | 説明 |
| --- | --- | --- | --- |
| POST | `/api/auth/signup` | 不要 | ユーザー登録 → JWT を body で返す |
| POST | `/api/auth/login` | 不要 | ログイン → JWT を body で返す (`{token, user}`) |
| POST | `/api/auth/logout` | 必要 | ログアウト (204。トークン破棄はクライアント側) |
| GET | `/api/auth/me` | 必要 | ログイン中のユーザー情報 |

### Todo

| メソッド | パス | 認証 | 説明 |
| --- | --- | --- | --- |
| GET | `/api/todos` | 必要 | 自分の Todo 一覧 |
| POST | `/api/todos` | 必要 | Todo 作成 (201 で作成した Todo を返す) |
| GET | `/api/todos/{id}` | 必要 | Todo 1 件取得 (他人の Todo は 404) |
| PUT | `/api/todos/{id}` | 必要 | Todo 更新 (部分更新。未指定フィールドは変更しない) |
| DELETE | `/api/todos/{id}` | 必要 | Todo 削除 (ソフトデリート、204) |

- **エラー形式**: `{ "code": string, "message": string }`
- **Todo は userID で絞り込み**。他人の Todo は取得・更新・削除できない (404)

### 認証フロー

```
1. signup / login → レスポンス body の token (JWT) を受け取る
2. クライアントは token を localStorage に保存
3. 認証必須 API には Authorization: Bearer <token> ヘッダを付与
4. logout はクライアントが localStorage から token を削除 (サーバは 204 を返すだけ)
```

リクエストは **2 段の middleware** を通ってから handler に届く：

```
リクエスト → [validation] OpenAPI spec と照合 (必須/形式 → 400)
          → [auth]       Bearer JWT 検証 (認証必須ルートのみ → 401)
          → handler
```

- どのルートが認証必須かは `api/openapi.yaml` の `security: bearerAuth` で宣言 → 生成コードがそのルートにのみ認証を要求する（spec 駆動）
- パスワードは bcrypt でハッシュ化して保存。JWT は `JWT_SECRET` で HS256 署名

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
| `http://localhost:8080/swagger/` | 📘 **API 仕様書（Swagger UI）** — OpenAPI 仕様をブラウザで閲覧（**ローカル限定**） |

`http://localhost:3000` でアクセスすると、フロントから `/api/hello` を叩いてもNext.js自身が応答してしまい、Goバックエンドに到達しません。Nginxを経由することでフロント・APIが同一オリジンに統一され、CORSも回避できます。

`http://localhost:8080/swagger/` は backend を直接叩く形で、OpenAPI 仕様 (`api/openapi.yaml`) をインタラクティブに閲覧・試打できます。本番（Cloud Run）では `APP_ENV=production` により自動的に無効化されます。

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
| `make mod-tidy` | `go mod tidy`（依存パッケージの整理） |

### データベース

| コマンド | 内容 |
| --- | --- |
| `make migrate-up` | 未適用のマイグレーションを全て適用 |
| `make migrate-down` | 直近のマイグレーションを 1 つロールバック |
| `make migrate-create name=<name>` | 新しいマイグレーションファイル（up/down）を `backend/migrations/` に生成 |
| `make migrate-version` | 現在のマイグレーションバージョンを表示 |

> golang-migrate を Docker イメージ (`migrate/migrate`) で実行する構成です。
> 接続文字列は [docker-compose.yml](docker-compose.yml) の `migrate` サービスの `entrypoint` で固定しているため、Makefile 側は短いサブコマンドだけになっています。

#### マイグレーションファイル作成の例

```bash
# 1. 空のマイグレーションファイル（up/down 2 つ）を生成
make migrate-create name=create_users
# → backend/migrations/000001_create_users.up.sql / .down.sql が生成される

# 2. 中身を SQL で記述（CREATE TABLE / DROP TABLE など）

# 3. DB に適用
make migrate-up

# 4. 現在のバージョンを確認
make migrate-version
```

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

### OpenAPI 自動生成

| コマンド | 内容 |
| --- | --- |
| `make gen-api` | Go と TS の型をまとめて再生成 |
| `make gen-api-backend` | Go の型のみ再生成（`backend/gen/api/api.gen.go`） |
| `make gen-api-frontend` | TS の型のみ再生成（`frontend/types/api.ts`） |

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
