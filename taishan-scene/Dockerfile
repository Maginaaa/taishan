FROM golang:1.20-alpine AS builder
WORKDIR /app

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

COPY . .
RUN go mod tidy && go mod download
# RUN go run cmd/generate/generate.go
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o target/main main.go


FROM golang:1.20-alpine
RUN apk add --no-cache tzdata
ENV TZ Asia/Shanghai

WORKDIR /app

COPY --from=builder /app/target/main /app/main

CMD ["/app/main"]
