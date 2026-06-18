# 1.ビルドステージ: 本番用バイナリを生成。コンパイラなどは本番には持っていかない
FROM golang:1.25 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server .

# 2.開発ステージ: docker-compose用 (ホットリロード)
FROM golang:1.25 AS dev

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

# ホットリロードを有効にするためのパッケージをインストール
RUN go install github.com/air-verse/air@v1.65.1

# ローカルでも make lint を動かせるよう golangci-lint を導入 (CI と同じ v2.11.4)
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh \
    | sh -s -- -b "$(go env GOPATH)/bin" v2.11.4

COPY . .

EXPOSE 8080

CMD ["air"]

# 3.本番ステージ: Cloud Run用 (distrolessでバイナリのみ軽量化 + ホットリロードを入れない)
FROM gcr.io/distroless/base-debian12:nonroot AS prod

WORKDIR /app

COPY --from=builder /out/server /app/server

EXPOSE 8080

USER nonroot:nonroot

CMD ["/app/server"]
