FROM golang:1.23-alpine AS build

ARG TARGETARCH
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=${TARGETARCH:-amd64}

WORKDIR /app
COPY go.mod ./
RUN go mod download
COPY . .
RUN go mod tidy
RUN go build -o nebulogo .

FROM alpine:latest
WORKDIR /app
COPY --from=build /app/nebulogo /app/nebulogo
COPY /app/template /app/template
EXPOSE 80
ENTRYPOINT ["/app/nebulogo"]
# Build command: docker build --build-arg TARGETARCH=$(uname -m) -t nebulogo .