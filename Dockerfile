
    FROM golang:1.24-alpine AS builder

    WORKDIR /src
    COPY go.mod go.sum ./
    RUN go mod download
    
    COPY . .
    RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o mail-service
    

    FROM scratch
    
    WORKDIR /app
    COPY --from=builder /src/mail-service .
    COPY --from=builder /src/db/migrations ./db/migrations

    USER 65532:65532

    ENTRYPOINT ["./mail-service"]
    EXPOSE 8080
    