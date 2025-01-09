#!/bin/bash

# Argument validation check
if [ "$#" -ne 1 ]; then
    echo "Usage: $0 <action>"
    exit 1
fi

function build() {
    go build -o bookings cmd/web/*.go
}

function run() {
    ./bookings
}

function test() {
    go test -v ./...
}

function coverage() {
    go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
}

function cleanup() {
    rm -rf bookings
}

function help() {
    echo "Available actions:"
    echo "  build     : Build the application"
    echo "  run       : Run the application"
    echo "  test      : Test the application"
    echo "  coverage  : Generate code coverage report"
    echo "  clean     : Clean up the application"
    echo "  help      : Display this help message"
}

case "$1" in
    "build")
        build
        ;;
    "run")
        build && run
        ;;
    "test")
        test
        ;;
    "coverage")
        coverage
        ;;
    "clean")
        cleanup
        ;;
    *)
        help
        ;;
esac