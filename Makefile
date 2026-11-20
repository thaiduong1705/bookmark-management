run:
	go run cmd/api/main.go

build:
	go build -o /dev/null cmd/api/main.go

swagger:
	swag init -g cmd/api/main.go

COVERAGE_EXCLUDE=mocks|main.go|test|docs|config.go|client.go
COVERAGE_THRESHOLD=100

test:
	go test ./... -coverprofile=coverage.tmp -covermode=atomic -coverpkg=./... -p 1
	grep -vE "$(COVERAGE_EXCLUDE)" coverage.tmp > coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@total=$$(go tool cover -func=coverage.out | grep total: | awk '{print $$3}' | sed 's/%//'); \
	if [ $$(echo "$$total < $(COVERAGE_THRESHOLD)" | bc -l) -eq 1 ]; then \
		echo "Coverage ($$total%) is below threshold ($(COVERAGE_THRESHOLD)%)"; \
		exit 1; \
	else \
		echo "Coverage ($$total%) meets threshold ($(COVERAGE_THRESHOLD)%)"; \
	fi

IMAGE = "$$DOCKER_USERNAME"
VERSION ?= latest

docker-test:
	docker build --target test --no-cache --progress=plain .

docker-build:
	docker build --target runner -t $(IMAGE)/bookmark_service:$(VERSION) .

docker-login:
	echo "$$DOCKER_PASSWORD" | docker login -u "$$DOCKER_USERNAME" --password-stdin

docker-push:
	docker push $(IMAGE):$(VERSION)
