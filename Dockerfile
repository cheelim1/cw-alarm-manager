FROM golang:1.21-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /bin/cw-alarm-manager main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN adduser -D appuser
USER appuser

WORKDIR /app

COPY --from=builder /bin/cw-alarm-manager /app/

CMD ["./cw-alarm-manager"]
