# =========================================
# build stage
# =========================================
FROM golang:1.17-alpine as build-stage

WORKDIR /opt/todo

ARG TZ
RUN apk --no-cache add tzdata && \
  cp /usr/share/zoneinfo/${TZ} /etc/localtime && \
  echo "${TZ}" > /etc/timezone && \
  apk del tzdata

COPY go.mod go.sum ./
RUN go mod download \
  && go install github.com/rakyll/gotest@latest \
  && go install github.com/cespare/reflex@latest \
  && go install github.com/golang/mock/mockgen@v1.6.0 \
  && go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
  && apk --no-cache add make alpine-sdk mysql-client

COPY ./ ./
COPY ./templates /opt/templates
RUN CGO_ENABLED=0 go build -o server

CMD ["/bin/sh"]

# =========================================
# debug stage
# =========================================
FROM golang:1.17-alpine as debug-stage

WORKDIR /opt/todo

ARG TZ
RUN apk --no-cache add tzdata && \
  cp /usr/share/zoneinfo/${TZ} /etc/localtime && \
  echo "${TZ}" > /etc/timezone && \
  apk del tzdata

COPY go.mod go.sum ./
RUN go mod download \
  && go install github.com/rakyll/gotest@latest \
  && go install github.com/cespare/reflex@latest \
  && go install github.com/golang/mock/mockgen@v1.6.0 \
  && go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@latest \
  && apk --no-cache add make alpine-sdk mysql-client \
  && CGO_ENABLED=0 go install github.com/go-delve/delve/cmd/dlv@latest

COPY ./ ./
COPY ./templates /opt/templates
RUN go build -gcflags "all=-N -l" -o debug_app

EXPOSE 8080 2345
CMD ["dlv", "--listen=:2345", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "./debug_app"]

# =========================================
# production stage
# =========================================
FROM alpine as production-stage

WORKDIR /opt/todo

ARG TZ
RUN apk --no-cache add tzdata && \
  cp /usr/share/zoneinfo/${TZ} /etc/localtime && \
  echo "${TZ}" > /etc/timezone && \
  apk del tzdata

COPY --from=build-stage /opt/todo/server /usr/bin/server
COPY ./templates /opt/templates

EXPOSE 8080

ENTRYPOINT ["server"]
