FROM golang:1.22-bookworm AS builder
WORKDIR /opt/workflow
COPY go.mod /opt/workflow
COPY go.sum /opt/workflow
COPY cmd /opt/workflow/cmd
COPY internal /opt/workflow/internal

ENV GOPROXY https://goproxy.cn,direct
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o /opt/workflow/server /opt/workflow/cmd/server/main.go

FROM alpine AS api-server
RUN apk update --no-cache && apk add --no-cache tzdata
ENV TZ Asia/Shanghai
WORKDIR /opt/workflow
COPY --from=builder /opt/workflow/server /opt/workflow/server
COPY ./config/config_docker.yaml /opt/workflow/config.yaml
ARG PORT=9090
EXPOSE $PORT
CMD ./server --config config.yaml