# たまたまコマンドと同じファイルがあると動かなくなることを防止
.PHONY: lint lint-fix up down rebuild-backend rebuild-frontend rebuild-nginx test frontend-install reset-frontend lint-frontend format-frontend format-check-frontend

# 1.静的解析を実行する
lint:
	docker compose exec backend golangci-lint run ./...
# 2.静的解析+自動修正を実行する
lint-fix:
	docker compose exec backend golangci-lint run --fix ./...
# 3.Dockerコンテナを起動する
up:
	docker compose up
# 4.Dockerコンテナを停止する
down:
	docker compose down
# 5.バックエンドをビルドし直して起動する
rebuild-backend:
	docker compose up -d --build backend
# 6.フロントエンドをビルドし直して起動する
rebuild-frontend:
	docker compose up -d --build frontend
# 7.Nginxをビルドし直して起動する
rebuild-nginx:
	docker compose up -d --build nginx

# 8.テストを実行する
test:
	docker compose exec backend go test ./...

# 9.フロントエンドの依存パッケージをインストールする（引数: pkg="prettier eslint-config-prettier")
frontend-install:
	docker compose exec frontend npm install --save-dev $(pkg)

# 10.フロントエンドを匿名ボリュームごと作り直す（node_modulesが壊れた時用）
reset-frontend:
	docker compose rm -fsv frontend
	docker compose up -d --build frontend

# 11.フロントエンドのESLintを実行する
lint-frontend:
	docker compose exec frontend npm run lint

# 12.フロントエンドのESLintを自動修正する
lint-fix-frontend:
	docker compose exec frontend npm run lint:fix

# 13.フロントエンドのTypeScriptの型チェックを実行する
typecheck-frontend:
	docker compose exec frontend npm run typecheck

# 14.フロントエンドにPrettierを一括適用する（ファイル書き換えあり）
format-frontend:
	docker compose exec frontend npm run format

# 15.フロントエンドがPrettierルール通りに整形されているかチェックする（CI想定、書き換えなし）
format-check-frontend:
	docker compose exec frontend npm run format:check

