GOBIN := $(shell go env GOPATH)/bin
export PATH := $(PATH):$(GOBIN)

.PHONY: user

user:
	buf generate user.proto --template '{"version":"v2","plugins":[{"local":"protoc-gen-go","out":"UserService","opt":"paths=source_relative"},{"local":"protoc-gen-go-grpc","out":"UserService","opt":"paths=source_relative"}]}'
