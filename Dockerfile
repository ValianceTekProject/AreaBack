FROM golang:1.25-alpine3.22 AS base
WORKDIR /app

FROM base AS deps
COPY go.mod go.sum ./
RUN go mod download

FROM base as builder
COPY --from=deps /go/pkg /go/pkg
COPY . .
RUN go run github.com/steebchen/prisma-client-go generate
RUN go build -o back .

FROM base AS final
COPY --from=builder /app/ .
EXPOSE 8080
CMD ["sh", "-c", "go run github.com/steebchen/prisma-client-go db push --accept-data-loss && go run seed/main.go && ./back"]
