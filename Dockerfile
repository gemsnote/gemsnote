# Stage 1: Build the Vue frontend
FROM node:22-alpine AS frontend-builder

WORKDIR /build/frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

# Stage 2: Build the Go application
FROM golang:1.24-alpine AS builder

# Set Go proxy for China
ENV GOPROXY=https://goproxy.cn,direct

# Install build dependencies
RUN apk add --no-cache git gcc musl-dev

WORKDIR /build

# Copy go mod files first for better caching
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .
COPY --from=frontend-builder /build/frontend/dist ./frontend/dist
# Never compile against a local deployment config. Docker Compose supplies its
# own config at runtime, while standalone images use the checked-in template.
RUN cp conf/app.conf-default conf/app.conf

# Generate the Revel application entrypoint, then build the server rather than
# the Revel command-line helper in app/cmd.
RUN go run ./app/cmd build . /tmp/gemsnote-revel-build prod && \
    CGO_ENABLED=0 GOOS=linux go build -trimpath -o /gemsnote -ldflags="-s -w" ./app/tmp

RUN REVEL_DIR="$(go list -m -f '{{.Dir}}' github.com/revel/revel)" && \
    mkdir -p /runtime/github.com/revel/revel && \
    cp -R "$REVEL_DIR/conf" "$REVEL_DIR/templates" /runtime/github.com/revel/revel/

# Stage 2: Runtime image
FROM alpine:3.22

LABEL maintainer="raptor<raptor.zh@gmail.com>"

# Install ca-certificates for HTTPS connections and timezone data
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /opt/gemsnote

# Copy the binary from builder
COPY --from=builder /gemsnote /opt/gemsnote/gemsnote
COPY --from=builder /runtime /opt/gemsnote/runtime

# Copy necessary runtime files
COPY conf/ /opt/gemsnote/conf/
COPY conf/app.conf-default /opt/gemsnote/conf/app.conf
COPY messages/ /opt/gemsnote/messages/
COPY public/ /opt/gemsnote/public/
COPY database/ /opt/gemsnote/database/
COPY --from=frontend-builder /build/frontend/dist/ /opt/gemsnote/frontend/dist/
COPY app/views/ /opt/gemsnote/app/views/
RUN mkdir -p /opt/gemsnote/files /opt/gemsnote/public/upload /opt/gemsnote/runtime/github.com/gemsnote && \
    ln -s /opt/gemsnote /opt/gemsnote/runtime/github.com/gemsnote/gemsnote

EXPOSE 9000

# Set working directory to where the app needs to run
WORKDIR /opt/gemsnote

CMD ["/opt/gemsnote/gemsnote", "-importPath", "github.com/gemsnote/gemsnote", "-srcPath", "/opt/gemsnote/runtime", "-runMode", "prod"]
