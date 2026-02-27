FROM golang:1.25-alpine AS builder

ARG VERSION=dev
ARG GIT_COMMIT=unknown
ARG BUILD_TIME=unknown

WORKDIR /src

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w \
    -X github.com/moabdelazem/gitops-sample-app/pkg/version.Version=${VERSION} \
    -X github.com/moabdelazem/gitops-sample-app/pkg/version.GitCommit=${GIT_COMMIT} \
    -X github.com/moabdelazem/gitops-sample-app/pkg/version.BuildTime=${BUILD_TIME}" \
    -o /bin/gitops-app ./cmd/main.go

FROM alpine:3.21

RUN apk --no-cache add ca-certificates \
    && addgroup -S appgroup \
    && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /bin/gitops-app .
COPY --from=builder /src/web ./web

USER appuser

EXPOSE 8080

ENTRYPOINT ["./gitops-app"]
