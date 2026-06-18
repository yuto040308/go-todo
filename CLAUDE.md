# go-todo

Go ポートフォリオプロジェクト。Phase 1 ゴール = CRUD + 認証認可 + CI/CD + Linter を備えた TODO アプリを Vercel / Cloud Run / Supabase にデプロイ

## 関連リンク

- Notion 親 PJ: https://www.notion.so/346c058a77ac80d7a704c7de80bd3a5e
- Notion やること看板: https://www.notion.so/2ecc058a77ac8082944ccfe9a199b0cd
- 開発思想 (MMDD): https://www.notion.so/20260530-AI-MMDD-370c058a77ac807aa2f0defa427eff48
- メモリ (好み・履歴): `~/.claude/projects/-Users-onishiyuto-Desktop-go-todo/memory/MEMORY.md`

## 技術スタック

バージョンの真の SoT は `backend/go.mod` と `frontend/package.json`。下表は概要。

**フロントエンド**: Next.js 16 (App Router) / React 19 / TypeScript 5 / Tailwind v4 / shadcn/ui / TanStack Query / openapi-typescript
**バックエンド**: Go 1.25 / Gin / GORM / golang-jwt v5 / oapi-codegen v2 (型/サーバ I/F 生成) / swaggest/swgui (ローカル Swagger UI)
**DB**: PostgreSQL 18 (本番は Supabase)
**インフラ**: Docker (ローカル) / Vercel (フロント) / Cloud Run (バック) / Supabase (DB)

## フォルダ構造

- `api/openapi.yaml` — **Single Source of Truth**。API 型はここから生成
- `backend/` — Go (Gin) API サーバ
  - `gen/api/` — oapi-codegen 生成物 (**編集禁止**)
  - `models/` — GORM の DB モデル (UUID ベース、`gorm.Model` 不使用)
  - `middleware/` — CORS, AuthHandler (Bearer JWT 検証)
  - `database/`, `handler/`, `usecase/`, `migrations/`
- `frontend/` — Next.js
  - `types/api.ts` — openapi-typescript 生成物 (**編集禁止**)
  - `api/` — axios ラッパ層
- `nginx/` — リバプロ (フロントとバックを同一オリジンに統一)
- `docs/` — 概念解説 md (clean-architecture / SOLID 等)。仕様ではなく学習資料 (`.gitignore` 済)
- `.github/workflows/` — CI (Reviewdog, PR-Agent, go test, frontend lint)

## 重要 workflow

- **yaml を変えたら必ず `make gen-api`** (Go と TS の型をまとめて再生成)
- **DB 変更**: `make migrate-create name=xxx` → SQL 書く → `make migrate-up`
- **静的解析**: `make lint` (Go) / `make lint-frontend` / `make typecheck-frontend` / `make format-check-frontend` / `make unused-check-frontend`
- **依存整理**: `make mod-tidy`
- **Swagger UI**: http://localhost:8080/swagger/ (ローカル限定、`APP_ENV=production` で無効化)
- 全 Makefile ターゲットの詳細は [README.md](README.md) を参照

## 設計規約 (確定済)

- **レイヤー分離**: API DTO (`backend/gen/api/`) と DB Model (`backend/models/`) は別物。境界で変換する
- **ID**: UUID (`string + format:uuid`)。`gorm.Model` は使わず手動でフィールド宣言、`BeforeCreate` で `uuid.New()`
- **認証**: `Authorization: Bearer <JWT>` 方式 (Cookie ではない、localStorage 保管前提)
- **JSON 命名**: snake_case (`user_name`, `created_at`)
- **エラー形式**: `{code: string, message: string}` (`{error}` ではない)
- **生成物の commit 方針**: コミットする (CI で diff 検出を効かせるため、`.gitignore` しない)

## 進め方の流儀 (詳細は memory 参照)

- **MMDD で進める** — 理解は人間が、速度は AI が。設計の主導権はユーザー
- **イテレーティブに進める** — 一気に複数機能を足さない。機能ごとにブランチ
- **新技術導入時** — Claude は先回り実装せず、概念解説 → ユーザー実装 → 添削
- **破壊的アクション** (git reset, force push, drop table, --no-verify 等) は事前確認
- **本番マイグレーション** — 起動時 migrate は避け、④ Supabase チケットで Cloud Run Jobs として構築予定

## チケット作業の起点

- 新チケット着手 → `/go-todo <notion-url>` (Skill が要件整理 + 計画立案)
- 進行中タスク再開 → `/go-todo` (引数なし、現状サマリ表示)
