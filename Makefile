start:
	@go run src/main.go
lint:
	@golangci-lint run
tests:
	@go test -v ./test/...
tests-%:
	@go test -v ./test/... -run=$(shell echo $* | sed 's/_/./g')
testsum:
	@cd test && gotestsum --format testname
swagger:
	@cd src && swag init
docker:
	@docker-compose up --build
docker-test:
	@docker-compose up -d && make tests
docker-down:
	@docker-compose down --rmi all --volumes --remove-orphans
docker-cache:
	@docker builder prune -f