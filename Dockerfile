FROM golang:1.22.5 AS base

FROM base AS deps
WORKDIR /app

COPY go.mod ./
RUN go mod download

FROM deps AS builder
WORKDIR /app

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /go-whyye

FROM builder AS runner
WORKDIR /app

ENV PORT=${PORT:-8080}
EXPOSE ${PORT:?}
CMD ["/go-whyye"]

