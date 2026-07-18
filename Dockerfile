ARG GO_VERSION=1.25.11
ARG ALPINE_VERSION=3.24

FROM --platform=$BUILDPLATFORM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS builder

ARG TARGETOS
ARG TARGETARCH

ENV GOPROXY=https://goproxy.cn,direct \
    GOSUMDB=sum.golang.org \
    CGO_ENABLED=0

WORKDIR /src

# 1) 仅拷贝模块清单，单独一层缓存依赖下载
COPY go.mod go.sum* /src

RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

# 2) 仅拷贝编译所需源码，避免 public/config 变更触发重编
COPY app/ /src/app/
COPY cmd/ /src/cmd/
COPY internal/ /src/internal/
COPY main.go /src/

RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build,id=thirteen-admin-gobuild-${TARGETOS}-${TARGETARCH} \
    GOOS=${TARGETOS} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags="-s -w" -o /out/thirteen-admin-server ./main.go

# 应用程序执行镜像构建
FROM yongyuanfan/golang-app:1.25.11.beta1 AS runtime

WORKDIR /app

COPY --from=builder --chmod=0755 /out/thirteen-admin-server /app/thirteen-admin-server

COPY ./.env.example /app/.env.example
COPY ./config/settings.container.yml /app/config/settings.yml
COPY ./public/ /app/public/

ENV APP_BIN=/app/thirteen-admin-server

EXPOSE 8080

CMD ["server", "-c", "/app/config/settings.yml"]
