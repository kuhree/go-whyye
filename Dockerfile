FROM golang:1.22.5-alpine AS base
ENV USER=appuser
ENV UID=10001

RUN apk update && apk upgrade \ 
        && apk add --no-cache git ca-certificates tzdata \
        && update-ca-certificates

# Add user
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --shell "/sbin/nologin" \
    --no-create-home \
    --uid "${UID}" \
    "${USER}"

FROM base AS deps
WORKDIR /app

COPY go.mod go.sum package.json ./
RUN go mod download && apk add --no-cache ca-certificates build-base

FROM deps AS builder
ENV GOOS=linux
ENV CGO_ENABLED=1
WORKDIR /app

COPY . .
RUN mkdir -p ./out ./out/state ./out/share ./out/bin && \
  go build -ldflags '-w -s -extldflags "-static"' -a -o ./out/bin/main main.go

FROM scratch AS runner
ENV USER=appuser                        
ENV PORT=${PORT:-8080}
ENV APP_ENV=${APP_ENV:-production}
WORKDIR /app

COPY --from=base /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=base /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=base /etc/passwd /etc/passwd   
COPY --from=base /etc/group /etc/group
COPY --from=builder /app ./

USER appuser
EXPOSE ${PORT:?}
CMD ["./out/bin/main"]
