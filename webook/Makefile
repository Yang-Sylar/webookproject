.PHONY: docker
docker:
	@rm webook || true
	@GOOS=linux GOARCH=arm go build -o webook .
	@docker rmi -f yzletter/webook:v0.0.1 .
	@docker build -t yzletter/webook:v0.0.1 .