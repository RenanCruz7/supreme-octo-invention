# Build stage
FROM golang:1.26.2-alpine AS builder

WORKDIR /app

# Instalar dependências necessárias
RUN apk add --no-cache git gcc musl-dev

# Copiar go.mod e go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copiar código
COPY . .

# Build da aplicação
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o awesomeProject .

# Runtime stage
FROM alpine:latest

WORKDIR /root/

# Instalar ca-certificates para HTTPS
RUN apk --no-cache add ca-certificates

# Copiar binary do build stage
COPY --from=builder /app/awesomeProject .

# Expor porta (padrão)
EXPOSE 8080

# Variáveis de ambiente padrão
ENV PORT=:8080 \
    DB_DRIVER=postgres

# Executar aplicação
CMD ["./awesomeProject"]

