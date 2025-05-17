FROM golang:1.20-alpine AS builder
WORKDIR /app

ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct

COPY . .
RUN go mod tidy && go mod download
RUN CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o target/main main.go

FROM golang:1.20-alpine

RUN echo "Asia/Shanghai" > /etc/timezone

WORKDIR /app
COPY --from=builder /app/target/main /app/main
COPY --from=builder /app/config/ /app/config/

CMD ["/app/main"]
