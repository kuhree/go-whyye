FROM golang:1.22.5 AS base

FROM base AS deps
WORKDIR /app

COPY go.mod go.sum package.json ./
RUN go mod download

FROM deps AS builder
WORKDIR /app

COPY . .

ENV CGO_ENABLED=1
RUN mkdir -p ./out ./out/state ./out/share ./out/bin 
RUN go build -o ./out/bin/main main.go

FROM deps AS runner
WORKDIR /app

ENV PORT=${PORT:-8080}
ENV APP_ENV=${APP_ENV:-production}
EXPOSE ${PORT:?}

COPY --from=builder /app/out ./out
COPY --from=builder /app/static ./static 
COPY --from=builder /app/templates ./templates 
COPY --from=builder /app/pkg ./pkg
CMD ["./out/bin/main"]

