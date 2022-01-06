FROM golang:1.17-alpine as build-stage

RUN apk add --update-cache alpine-sdk
WORKDIR /opt

COPY go.mod go.sum ./
RUN go mod download \
  && go get -u github.com/rakyll/gotest \
  && go get -u github.com/cespare/reflex

COPY ./ ./
RUN go build -o server

FROM golang:1.17-alpine as debug-stage

RUN apk add --update-cache alpine-sdk
WORKDIR /opt

RUN CGO_ENABLED=0 go install github.com/go-delve/delve/cmd/dlv@latest
COPY . /opt
RUN go build -gcflags "all=-N -l" -o debug_app
EXPOSE 8080 2345
ENTRYPOINT ["dlv", "--listen=:2345", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "./debug_app"]

FROM alpine
COPY --from=build-stage /opt/server /usr/bin/server
EXPOSE 8080
ENTRYPOINT ["server"]
