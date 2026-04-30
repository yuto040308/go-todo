# 1.ベースイメージの指定
FROM golang:1.26.2

# 2.コンテナ上にワークディレクトリを作成
WORKDIR /app

# 3.コンテナにgo.modとgo.sumをコピーして、必要なライブラリをインストール
COPY go.mod go.sum ./
RUN go mod download

# 4.ホットリロードを有効にするためのパッケージをインストール
RUN go install github.com/air-verse/air@v1.65.1

# 5.静的コード分析をインストール
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
    | sh -s -- -b $(go env GOPATH)/bin v2.11.4

# 6.ホストのbackend/配下のファイルをコンテナにコピー
COPY . .

# 7.ポートの公開
EXPOSE 8080

# 8.Airを起動してソースの変更を監視しながらgo build -> 再起動をループさせる
CMD ["air"]
