FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.sum ./

COPY go.mod ./

RUN go mod download

COPY *.go ./

RUN CGO_ENABLED=0 GOOS=linux go build -o weather-api .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/weather-api .

COPY --from=builder /app/go.mod .

USER appuser 

EXPOSE 8000

CMD ["./weather-api"]
