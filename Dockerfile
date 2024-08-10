ARG APP_ENV="production"

FROM golang:1.22.5-alpine AS deps
  ENV GOOS=linux
  ENV CGO_ENABLED=1
  WORKDIR /app

  COPY go.mod go.sum package.json ./
  RUN apk update && apk upgrade \ 
          && apk add --no-cache git ca-certificates tzdata build-base \
          && update-ca-certificates \
          && go mod download

FROM deps AS builder
  WORKDIR /app

  COPY . .
  RUN mkdir -p ./out ./out/state ./out/share ./out/bin \
    && go build -ldflags '-w -s -extldflags "-static"' -a -o ./entrypoint main.go

FROM getsentry/sentry-cli:latest AS release
  ENV APP_ENV=${APP_ENV}
  ARG SKIP_RELEASE="1"
  ENV SKIP_RELEASE=${SKIP_RELEASE}
  ARG SENTRY_RELEASE_FALLBACK="docker"
  ENV SENTRY_RELEASE_FALLBACK=${SENTRY_RELEASE_FALLBACK}
  ARG SENTRY_AUTH_TOKEN=""
  ENV SENTRY_AUTH_TOKEN=${SENTRY_AUTH_TOKEN}
  WORKDIR /app

  COPY --from=builder /app ./
  RUN ./bin/sentry-release.sh

FROM alpine:latest AS runner
  ENV USER=appuser
  ENV UID=10001
  ENV PORT=${PORT:-8080}
  ENV APP_ENV=${APP_ENV}
  WORKDIR /app

  COPY --from=deps /usr/share/zoneinfo /usr/share/zoneinfo
  COPY --from=deps /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
  COPY --from=deps /etc/passwd /etc/passwd   
  COPY --from=deps /etc/group /etc/group
  COPY --from=release /app ./

  # Add user
  RUN adduser \
      --disabled-password --gecos "" --home "/nonexistent" --shell "/sbin/nologin" \
      --no-create-home --uid "${UID}" "${USER}" \
      && chown -R ${USER}:${USER} /app 

  USER appuser
  EXPOSE ${PORT:?}
  CMD ["./entrypoint"]
