install_go_dependencies() {
    echo "Installing go dependencies..."
    go get -tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
    go get -tool  go.uber.org/mock/mockgen@latest
    go get -tool github.com/securego/gosec/v2/cmd/gosec@latest
    go get -tool github.com/google/wire/cmd/wire@latest
    go get -tool go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    go get -tool go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
}

install_npm_dependencies() {
    npm install -g openapi-format@1.16.0
    npm install -g @redocly/cli@1.9.1
}

install_npm_dependencies
install_go_dependencies