# 1.ベースイメージの指定
FROM golang:1.26.2

# 2.コンテナ上にワークディレクトリを作成
WORKDIR /app

# 3.コンテナにgo.modとgo.sumをコピーして、必要なライブラリをインストール
COPY go.mod go.sum ./
RUN go mod download

# 4.ホットリロードを有効にするためのパッケージをインストール
RUN go install github.com/air-verse/air@v1.65.1

# 5.ホストのbackend/配下のファイルをコンテナにコピー
COPY . .

# 6.ポートの公開
EXPOSE 8080

# 7.Airを起動してソースの変更を監視しながらgo build -> 再起動をループさせる
CMD ["air"]
