#!/bin/bash

go test -v -coverprofile=size_coverage.out ./... && go  tool cover -html ./size_coverage.out -o ./cover.html
#go test -v -coverprofile=size_coverage.out ./encoder/... && go  tool cover -html ./size_coverage.out -o ./cover.html
#go test -v -coverprofile=size_coverage.out ./open_standard/... && go  tool cover -html ./size_coverage.out -o ./cover.html
#go test -v -coverprofile=size_coverage.out ./runtime/... && go  tool cover -html ./size_coverage.out -o ./cover.html
#go test -v -coverprofile=size_coverage.out ./field/... && go  tool cover -html ./size_coverage.out -o ./cover.html
#go test -v -coverprofile=size_coverage.out ./log/... && go  tool cover -html ./size_coverage.out -o ./cover.html
