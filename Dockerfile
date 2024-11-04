LABEL authors="nhpro"

FROM golang:1.20-alpine AS build

ARG TARGETARCH
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=${TARGETARCH:-amd64}

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o nebulogo .

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/nebulogo /app/nebulogo
EXPOSE 80
ENTRYPOINT ["/app/nebulogo"]
# Build command: docker build --build-arg TARGETARCH=$(uname -m) -t nebulogo .