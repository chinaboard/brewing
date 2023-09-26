
build-docker-runner:
	docker build -t brewing-runner -f runner.Dockerfile .

build-cli:
	bash scripts/build-cli.sh

build-docker-service:
	docker build -t brewing-service -f service.Dockerfile .

build-service:
	bash scripts/build-service.sh
