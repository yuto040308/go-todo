# 1.ベースイメージの指定
FROM golang:1.26.2

# 2.コンテナにgo.modとgo.sumをコピーして、必要なライブラリをインストール
COPY go.mod go.sum ./
RUN go mod download

# 3.コンテナ上にワークディレクトリを作成
WORKDIR /app

# 4.ホストのbackend/配下のファイルをコンテナにコピー
COPY . .

# 5.アプリケーションのビルド
RUN go build -o main main.go

# 6.ポートの公開
EXPOSE 8080

# 7.コンテナ起動時に「/app/main」（実態ファイル）が実行されるように設定
CMD ["./main"]