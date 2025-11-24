FROM golang:1.25-alpine3.22 AS base
WORKDIR /app

FROM base AS deps
COPY go.mod go.sum ./
RUN go mod download

FROM base as builder
COPY --from=deps /go/pkg /go/pkg
COPY . .
RUN go build -o back .

FROM alpine:3.21 AS final
COPY --from=builder /app/back .
EXPOSE 8080
CMD ["./back"]
