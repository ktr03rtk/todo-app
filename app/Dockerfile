FROM golang as build-stage

WORKDIR /opt
COPY go.mod go.sum ./
RUN go mod download \
  && go get -u github.com/rakyll/gotest \
  && go get -u github.com/cespare/reflex

COPY ./ ./
RUN go build -o server

FROM golang as debug-stage
RUN go install github.com/go-delve/delve/cmd/dlv@latest
WORKDIR /opt
COPY . /opt
RUN go build -gcflags "all=-N -l" -o debug_app
EXPOSE 8080 2345
ENTRYPOINT ["dlv", "--listen=:2345", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "./debug_app"]

FROM golang
COPY --from=build-stage /opt/todo-app /usr/bin/todo-app
EXPOSE 8080
ENTRYPOINT ["todo-app"]
