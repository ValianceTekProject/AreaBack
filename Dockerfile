FROM golang:1.25-alpine3.22 AS builder

WORKDIR /app

# to update later
# COPY go.mod go.sum ./
COPY go.mod ./

RUN go mod download

COPY . .

RUN go build -o back .

FROM alpine:3.21 AS final

COPY --from=builder /app/back .

EXPOSE 8080

CMD ["./back"]
