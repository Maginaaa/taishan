#!/bin/bash

# 动态获取当前脚本所在目录作为工作区根目录
WORKSPACE_ROOT="$(cd "$(dirname "$0")" && pwd)"
TAG=1.0.1

# 定义需要构建的Go项目列表
GO_PROJECTS=(
  "taishan-account"
  "taishan-collector"
  "taishan-data"
  "taishan-machine"
  "taishan-report"
  "taishan-scene"
  "taishan-task"
  "taishan-engine"
)

# 构建每个Go项目
for PROJECT in "${GO_PROJECTS[@]}"
do
  echo "========================================"
  echo "正在构建项目: $PROJECT"
  echo "========================================"

  # 进入项目目录
  cd "$WORKSPACE_ROOT/$PROJECT" || exit

  # 执行go build
  echo "执行 go build..."
  CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -o "$PROJECT" main.go

  # 执行docker build
  echo "执行 docker build..."
  docker build --platform linux/amd64 -t "registry.cn-guangzhou.aliyuncs.com/your-repo/$PROJECT:$TAG" .
done

echo "fe build"
cd "$WORKSPACE_ROOT/taishan-fe" || exit
pnpm run build:pro
docker build --platform linux/amd64 -t "registry.cn-guangzhou.aliyuncs.com/your-repo/taishan-fe:$TAG" .

echo "所有项目构建完成！"


echo "done"
