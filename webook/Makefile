.PHONY: docker
docker:
	@rm webook || true
	@GOOS=linux GOARCH=arm go build -tags=k8s -o webook .
	@docker rmi -f yzletter/webook:v0.0.1 .
	@docker build -t yzletter/webook:v0.0.1 .

.PHONY: mock
mock:
	@mockgen -source=internal/service/code.go -package=svcmocks -destination=internal/service/mock/code.mock.go
	@mockgen -source=internal/service/user.go -package=svcmocks -destination=internal/service/mock/user.mock.go
	@go mod tidy
