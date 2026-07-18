# clash-admin 构建与部署 Makefile
#
# 常用命令：
#   make proto                          重新生成 gRPC 代码
#   make build-linux                    交叉编译 Linux amd64 二进制
#   make push-test                      上传二进制到测试服务器
#   make deploy-test                    编译并上传二进制到测试环境
#   make build-image                    构建多架构 Docker 镜像（不推送）
#   make push-image                     构建并推送 Docker 镜像到测试私有仓库
#   make push-image IMAGE_TAG=1.0.4     指定镜像版本标签

# ---------------------------------------------------------------------------
# 本地二进制构建
# ---------------------------------------------------------------------------

# 输出目录与可执行文件名
BIN_DIR := bin
BIN := $(BIN_DIR)/clash-admin-server

# 交叉编译目标平台（测试服务器为 Linux amd64）
GOOS := linux
GOARCH := amd64
# -w 去除 DWARF 调试信息，-s 去除符号表，减小二进制体积
LDFLAGS := -w -s

# ---------------------------------------------------------------------------
# 测试服务器（二进制 scp 上传）
# ---------------------------------------------------------------------------

REMOTE_USER := huyong
REMOTE_HOST := 148.70.103.64
REMOTE_PATH := ~/

# ---------------------------------------------------------------------------
# Docker 镜像（测试环境私有仓库）
# ---------------------------------------------------------------------------

# 测试环境 registry 地址
DOCKER_REGISTRY := 148.70.103.64:5000
DOCKER_IMAGE := clash-admin-server
# 镜像标签，可通过命令行覆盖，例如：make push-image IMAGE_TAG=1.0.4
IMAGE_TAG ?= latest
# 多架构构建目标：Apple Silicon / ARM 服务器 + x86_64 服务器
# DOCKER_PLATFORMS := linux/arm64,linux/amd64
DOCKER_PLATFORMS := linux/amd64
# buildx 跨次构建缓存 tag，与业务镜像 tag 独立
CACHE_REF := $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):buildcache

# 交叉编译 Linux amd64 静态二进制，产物位于 bin/clash-admin-server
build-linux:
	mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) \
		go build -ldflags="$(LDFLAGS)" -a -installsuffix cgo -o $(BIN) .

# 将已编译的二进制 scp 上传到测试服务器
# 前置条件：本机已配置到 REMOTE_HOST 的 SSH 免密或密钥登录
# 注意：需先执行 make build-linux，或直接使用 make deploy-test
push-test:
	scp $(BIN) $(REMOTE_USER)@$(REMOTE_HOST):$(REMOTE_PATH)

# 编译 Linux 二进制并上传到测试服务器（push-test 的完整流程）
deploy-test: build-linux push-test

# 使用 buildx 构建多架构 Docker 镜像，仅写入 buildx 缓存，不推送到 registry
# 适用于本地验证 Dockerfile 或构建流程
build-image:
	docker buildx build \
		--platform $(DOCKER_PLATFORMS) \
		--cache-from type=registry,ref=$(CACHE_REF) \
		--cache-to type=registry,ref=$(CACHE_REF),mode=max \
		-t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(IMAGE_TAG) \
		.

# 构建多架构镜像并推送到测试私有仓库
# 前置条件：已 docker login $(DOCKER_REGISTRY)
# 多平台 manifest 需 --push 才能完整发布，无法仅用 docker push 分步推送
push-image:
	docker buildx build \
		--platform $(DOCKER_PLATFORMS) \
		--cache-from type=registry,ref=$(CACHE_REF) \
		--cache-to type=registry,ref=$(CACHE_REF),mode=max \
		-t $(DOCKER_REGISTRY)/$(DOCKER_IMAGE):$(IMAGE_TAG) \
		--push \
		.
