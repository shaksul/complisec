FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install dependencies
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Final stage
FROM alpine:latest

RUN apk --no-cache add \
    ca-certificates \
    tesseract-ocr \
    poppler-utils \
    wget

# Install specific language data (rus, kaz, eng) into default tessdata path
RUN mkdir -p /usr/share/tessdata && \
    wget -qO /usr/share/tessdata/rus.traineddata https://github.com/tesseract-ocr/tessdata_fast/raw/main/rus.traineddata && \
    wget -qO /usr/share/tessdata/kaz.traineddata https://github.com/tesseract-ocr/tessdata_fast/raw/main/kaz.traineddata && \
    wget -qO /usr/share/tessdata/eng.traineddata https://github.com/tesseract-ocr/tessdata_fast/raw/main/eng.traineddata
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Expose port
EXPOSE 8080

# Run the binary
CMD ["./main"]
