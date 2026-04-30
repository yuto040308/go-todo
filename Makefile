# たまたまコマンドと同じファイルがあると動かなくなることを防止
.PHONY: lint lint-fix up down rebuild-backend rebuild-frontend rebuild-nginx

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
