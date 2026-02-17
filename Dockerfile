# builder 
FROM golang:1.25.5-alpine3.21 AS builder

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download && go mod verify

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build \
  -ldflags="-w -s" \
  -trimpath \
  -o weather-api \
  .
# run time 
FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=builder /etc/passwd /etc/passwd

COPY --from=builder /etc/group  /etc/group

WORKDIR /app

COPY --from=builder /app/weather-api .

USER nobody:nobody

EXPOSE 8000

ENTRYPOINT ["./weather-api"]
