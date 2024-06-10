FROM golang:1.22.1-alpine3.19 AS BuildStage
RUN apk update && apk add git
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN CGO_ENABLED=0 go build -ldflags "-X main.version=`git describe --tags --abbrev=0`" -o /app/memdods ./main.go

FROM alpine:latest
WORKDIR /
COPY --from=BuildStage /app/medods medods
COPY --from=BuildStage /app/config /config
ENTRYPOINT ["./medods"]