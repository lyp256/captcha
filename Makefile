gen:
	@echo generate...
	@go list ./... |grep -v vendor |xargs go generate

vet:
	@echo 静态语法检查...
	@go list ./... |grep -v vendor |xargs go vet

fmt:
	@echo 代码格式化...
	@go list ./... |grep -v vendor |xargs go fmt

lint:
	@echo lint...
	@go list ./... |grep -v vendor |xargs golint

tidy:
	@echo 依赖检查...
	@go mod tidy

vendor:
	@echo vendor
	@go mod vendor

test:
	@echo 单元测试...
	@go list ./... |grep -v vendor |xargs go test --cover -v

pretty: tidy fmt lint vet vendor

.PHONY: tidy fmt lint vet test vendor pretty
