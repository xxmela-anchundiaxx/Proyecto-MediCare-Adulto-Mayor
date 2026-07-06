# === ETAPA 1: BUILDER (Compilación) ===
# Usamos golang:latest para asegurarnos de tener la versión de Go (1.25+) que exige tu go.mod
FROM golang:latest AS builder

# Instalar git por si tus dependencias de Go lo requieren para descargarse
RUN apt-get update && apt-get install -y git

WORKDIR /app

# Copiar archivos de dependencias de Go primero
COPY go.mod go.sum ./
RUN go mod download

# Copiar todo el código fuente del proyecto
COPY . .

# Compilar apuntando al directorio de tu main.go de forma estática
RUN CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api/main.go

# === ETAPA 2: RUNNER (Imagen ligera final) ===
FROM alpine:latest AS runner

WORKDIR /app

# Instalar certificados e intérpretes mínimos para asegurar la ejecución del binario estático
RUN apk --no-cache add ca-certificates

# Copiar el binario compilado desde la etapa anterior
COPY --from=builder /app/main .

# Exponer el puerto de la API
EXPOSE 8080

# Comando para ejecutar la aplicación
CMD ["./main"]