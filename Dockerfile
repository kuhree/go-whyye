FROM golang:1.22.5-alpine AS base
  ENV GOOS=linux
  ENV CGO_ENABLED=1
  WORKDIR /app

  RUN apk update && apk upgrade \ 
          && apk add --no-cache git ca-certificates tzdata \
          && update-ca-certificates

  COPY go.mod go.sum package.json ./
  RUN go mod download && apk add --no-cache ca-certificates build-base

  COPY . .
  RUN mkdir -p ./out ./out/state ./out/share ./out/bin \
    && go build -ldflags '-w -s -extldflags "-static"' -a -o ./entrypoint main.go

FROM getsentry/sentry-cli:latest AS release
  WORKDIR /app
  ARG SKIP_RELEASE="1"
  ENV SKIP_RELEASE=${SKIP_RELEASE}
  ARG SENTRY_AUTH_TOKEN=""
  ENV SENTRY_AUTH_TOKEN=${SENTRY_AUTH_TOKEN}

  COPY --from=base /app ./
  RUN ./bin/sentry-release.sh

FROM alpine:latest AS runner
  ENV USER=appuser
  ENV UID=10001
  ENV PORT=${PORT:-8080}
  ENV APP_ENV=${APP_ENV:-production}
  WORKDIR /app

  COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
  COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
  COPY --from=base /etc/passwd /etc/passwd   
  COPY --from=base /etc/group /etc/group
  COPY --from=release /app ./

  # Add user
  RUN adduser \
      --disabled-password --gecos "" --home "/nonexistent" --shell "/sbin/nologin" \
      --no-create-home --uid "${UID}" "${USER}" \
      && chown -R ${USER}:${USER} /app 

  USER appuser
  EXPOSE ${PORT:?}
  CMD ["./entrypoint"]
