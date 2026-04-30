# go-todo

## 概要

Next.js + Go + PostgreSQL で構築した、シンプルな **ToDo 管理アプリ**です。

ToDo の CRUD（作成・一覧・更新・削除）をコア機能としつつ、**CI/CD パイプライン**や **AI コードレビュー**をはじめとするモダンな開発手法・最新技術を実践的に組み込むことを目的としたプロジェクトです。
「小さいアプリで、フルスタック開発のベストプラクティスをすべて体現する」ことをコンセプトに掲げています。

## 特徴

- シンプルな ToDo CRUD（作成・一覧・更新・削除）
- OpenAPI 仕様書をコードコメントから自動生成し、フロントエンドの型も自動生成
- GitHub Actions による CI/CD パイプライン（ビルド・テスト・デプロイの完全自動化）
- AI コードレビュー（Claude Code など）を開発フローに統合
- Vercel / Cloud Run / Supabase を活用したサーバーレス・マネージドな本番環境

---

## 技術構成

### フロントエンド

| カテゴリ | 技術 | バージョン | 用途 |
| --- | --- | --- | --- |
| フレームワーク | Next.js (App Router) | 15.x | Web フレームワーク |
| 言語 | TypeScript | 5.x | 型安全な JavaScript |
| スタイリング | Tailwind CSS | v4 | ユーティリティファーストの CSS |
| UI コンポーネント | shadcn/ui | - | Tailwind ベースのコンポーネント集 |
| データフェッチ | TanStack Query | v5.99.2 | サーバー状態管理・キャッシュ |
| 型自動生成 | openapi-typescript | v7.13.0 | OpenAPI 仕様書から TypeScript 型を生成 |

---

### バックエンド

| カテゴリ | 技術 | バージョン | 用途 |
| --- | --- | --- | --- |
| 言語 | Go | 1.26.2 | バックエンド言語 |
| フレームワーク | Gin | v1.12.0 | HTTP Web フレームワーク |
| API ドキュメント | swaggo/swag | v1.16.6 | コードコメントから OpenAPI 仕様書（yaml）を自動生成 |

> swaggo/swag は v2.0.0-rc5（RC版）もありますが、安定版の v1.16.6 を採用しています。

---

### データベース

| カテゴリ | 技術 | バージョン | 用途 |
| --- | --- | --- | --- |
| DB | PostgreSQL | 17.x | リレーショナルデータベース |

> Supabase が対応しているバージョンに合わせて 17.x を採用しています。

---

### インフラ

| サービス | 用途 |
| --- | --- |
| Vercel | Next.js フロントエンドのホスティング・デプロイ |
| Cloud Run | Go / Gin API サーバーのホスティング（Google Cloud） |
| Supabase | PostgreSQL のホスティング・認証・ストレージ |

---

## ローカル開発

### 起動

```bash
docker compose up -d --build
```

### アクセス

ブラウザからは **必ず Nginx 経由（ポート80）でアクセス**してください。

| URL | 用途 |
| --- | --- |
| `http://localhost` | ✅ Nginx経由（推奨）。`/api/*` がバックエンドに振り分けられる |
| `http://localhost:3000` | ❌ Next.js直アクセス。`/api/*` が解決できず 404 になる |

`http://localhost:3000` でアクセスすると、フロントから `/api/hello` を叩いてもNext.js自身が応答してしまい、Goバックエンドに到達しません。Nginxを経由することでフロント・APIが同一オリジンに統一され、CORSも回避できます。

### コード品質チェック（静的解析）

バックエンドコンテナには **golangci-lint** が同梱されています。Goコードの静的解析・コーディング規約チェックが実行できます。
具体的なコマンドは下記の「よく使うコマンド」を参照してください。

---

### よく使うコマンド（Makefile チートシート）

プロジェクトルートに `Makefile` を用意しているので、長いコマンドを覚えなくても `make xxx` で簡単に操作できます。

| コマンド | 内容 |
| --- | --- |
| `make up` | 全コンテナを起動 |
| `make down` | 全コンテナを停止 |
| `make rebuild-backend` | バックエンド（Go）をビルドし直して起動 |
| `make rebuild-frontend` | フロントエンド（Next.js）をビルドし直して起動 |
| `make rebuild-nginx` | Nginx をビルドし直して起動 |
| `make lint` | Go の静的解析を実行 |
| `make lint-fix` | Go の静的解析 + 自動修正 |

#### 使用例

```bash
# コンテナをまとめて起動
make up

# バックエンドのコードを変えた後、ビルドし直して起動
make rebuild-backend

# 静的解析を実行
make lint
```

> `./...` は「カレントディレクトリ以下のすべてのGoパッケージを対象」というGo標準の表記です。
> golangci-lint は `errcheck` `staticcheck` `govet` など複数のlinterを束ねたメタリンターで、デフォルト設定でも実用十分な指摘が得られます。

