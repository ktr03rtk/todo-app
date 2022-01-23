FROM golang:1.17-alpine as build-stage

WORKDIR /opt/todo

ARG TZ
RUN apk add -U tzdata && \
  cp /usr/share/zoneinfo/${TZ} /etc/localtime && \
  echo "${TZ}" > /etc/timezone && \
  apk del tzdata

COPY go.mod go.sum ./
RUN go mod download \
  && go get -u github.com/rakyll/gotest \
  && go get -u github.com/cespare/reflex \
  && apk -U add make alpine-sdk mysql-client

COPY ./ ./
RUN CGO_ENABLED=0 go build -o server

##################
FROM golang:1.17-alpine as debug-stage

WORKDIR /opt/todo
RUN CGO_ENABLED=0 go install github.com/go-delve/delve/cmd/dlv@latest
COPY . /opt/todo
RUN go build -gcflags "all=-N -l" -o debug_app

EXPOSE 8080 2345
ENTRYPOINT ["dlv", "--listen=:2345", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "./debug_app"]

##################
FROM alpine as production-stage
COPY --from=build-stage /opt/todo/server /usr/bin/server
EXPOSE 8080
ENTRYPOINT ["server"]
