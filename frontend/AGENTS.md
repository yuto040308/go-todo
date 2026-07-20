<!-- BEGIN:nextjs-agent-rules -->

# This is NOT the Next.js you know

This version has breaking changes — APIs, conventions, and file structure may all differ from your training data. Read the relevant guide in `node_modules/next/dist/docs/` before writing any code. Heed deprecation notices.

<!-- END:nextjs-agent-rules -->

---

# フロントエンド設計方針 (PF 1-4 で決定)

⑤モックを踏まえ、⑦実装の指針として決めた方針。**認証は Bearer + JWT / localStorage 保管**が前提（2026-05-25 に Cookie 方式から変更）。この前提が SC/CC・データ取得・ページ保護・フォームの全方針を規定している。

## ディレクトリ構成

```
frontend/
├─ app/
│  ├─ layout.tsx              # RootLayout (Providers 注入)
│  ├─ page.tsx                # トップ (→ /login or /todos へ誘導)
│  ├─ login/page.tsx          # 認証不要
│  ├─ signup/page.tsx         # 認証不要
│  └─ todos/                  # 認証必須。todos/layout.tsx が配下を一括ガード
│     ├─ layout.tsx           #   クライアント側 AuthGuard + 共通ヘッダー
│     ├─ page.tsx             #   一覧
│     ├─ new/page.tsx         #   新規作成
│     └─ [id]/edit/page.tsx   #   編集
├─ components/
│  ├─ shadcn/                 # shadcn/ui 由来の所有プリミティブ (下記②)
│  ├─ features/               # 機能特化 (todos/ auth/ でサブ分割、hooks も同居)
│  └─ layout/                 # ヘッダー等の共通レイアウト部品
├─ lib/
│  ├─ api/                    # API 層を集約
│  │  ├─ client.ts            #   axios インスタンス (interceptor)
│  │  ├─ todos.ts             #   Todo 系の呼び出し関数
│  │  └─ auth.ts              #   login / signup / me
│  ├─ auth-context.tsx        # AuthContext / useAuth
│  └─ utils.ts               # cn() 等
└─ types/
   └─ api.ts                  # openapi-typescript 生成 (編集禁止)
```

### コンポーネントの3カテゴリ（触ってよいかの境界）

| カテゴリ                | 実体                                                     | 編集              | 再生成                      |
| ----------------------- | -------------------------------------------------------- | ----------------- | --------------------------- |
| ① 編集禁止・自動再生成  | `types/api.ts`（openapi-typescript）、backend `gen/api/` | ❌                | `make gen-api` で毎回上書き |
| ② scaffold 由来だが所有 | `components/shadcn/`（shadcn/ui プリミティブ）           | ✅ 所有コード扱い | `shadcn add` は初回生成のみ |
| ③ 完全自作              | `components/features/`・`components/layout/`             | ✅                | 生成されない                |

> shadcn は「触るな」枠ではなく所有コード。`ui` ではなく `shadcn` と命名し、出自を明示しつつ編集可の含意を残す。配置先は `components.json` の `aliases.ui` で `@/components/shadcn` を指す。

## Server / Client Component の使い分け

**ハイブリッド**。localStorage の token は Server Component から読めないため、認証データを触る処理は Client 側になる。

- 静的な殻（レイアウト・見出し・カード枠）は **Server Component**（デフォルト）
- **データ取得 / mutation を含む「島」だけ** `'use client'` の Client Component に切り出す
- page はできるだけ薄く保ち、UI ロジックは `components/features/` に寄せる

## データ取得

全ケース **Client + TanStack Query**（`useQuery` / `useMutation`）で統一。`Providers`（`app/providers.tsx`）で `QueryClientProvider` を全体注入済み。Server Component での fetch は使わない（token を持てないため）。

## 認証状態管理

- **`AuthContext`（`lib/auth-context.tsx`）** … ログイン状態・`login()` / `logout()`・token の localStorage 出し入れを担う
- **`useQuery(['me'])`** … ログインユーザー情報はサーバ状態として React Query が管理（二重管理を避ける）
- token 本体は `lib/api/client.ts` の interceptor が localStorage から読んで `Authorization: Bearer` に載せる（実装済み）

> **⑦への申し送り**: ユーザー情報取得に `GET /me`（自分の情報を返す）API が必要。バックエンドに無ければ⑦で追加する。無い場合の簡略策は「login レスポンスに user を含めて Context に入れる」。

## ページ保護（認証必須ページ）

localStorage 方式のため `middleware.ts` では JWT を読めない → **クライアント側ガード**。

- `app/todos/layout.tsx`（Client Component）で `useAuth` を判定し、未認証なら `router.replace('/login')`
- 保護対象は今のところ `todos/` 配下のみ。`todos/layout.tsx` が配下 page（一覧 / new / [id]/edit）を一括で守る
- **判定が済むまでは `null` / スケルトンを返す**（保護ページが一瞬描画されるのを防ぐ）
- ⑦時点では **Route Group `(protected)/` は使わない判断**（保護ページが todos 配下に閉じているため過剰）。todos 外に認証必須ページが増えたら `(protected)/` グループ + 共有 layout に見直す

## フォーム & バリデーション

- 送信は **Client + `useMutation`**。`useActionState` / Server Action は使わない（サーバ側で token を扱えず認証方式と噛み合わないため撤回）
- バリデーションは **クライアント側の軽い検証**（`title` 必須・長さ等、素の state で十分）。ライブラリ（zod 等）は現状の要件には過剰
- **最終的な検証はバックエンド**（spec 駆動の validation middleware）が担う。フロントは UX 用と割り切る

## その他の確定事項

- **shadcn/ui**: 導入する（`components/shadcn/`、style は `radix-nova`、icon は lucide）
- **API ラッパ**: 軽い axios wrapper（`lib/api/client.ts`）+ 機能別関数（`lib/api/todos.ts` 等）。`fetch` 直書きはしない
- **JSON 命名**: snake_case（`user_name`, `created_at` 等、バックエンドと統一）

## ⑦（フロント実装）への移行タスク

このチケットはスケルトン（空フォルダ）と方針文書まで。以下は⑦で実施:

- `components/ui/` → `components/shadcn/` へリネーム（`components.json` の `aliases.ui` 更新 + 既存7コンポーネント移動 + import 修正）
- `lib/api.ts` → `lib/api/client.ts` へ移動、`api/hello.ts` の関数群を `lib/api/*.ts` へ集約
- `app/todos/layout.tsx`（AuthGuard）を追加し、`app/todos/page.tsx` のモックを実データ（React Query）に置換
- `GET /me` API の用意（または簡略策）
