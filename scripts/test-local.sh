#! /usr/bin/env bash
set -e

# run all tests with coverage
go test ./... -coverprofile=test/coverage/cover.out
# opens default browser with coverage report
go tool cover -html=test/coverage/cover.out
