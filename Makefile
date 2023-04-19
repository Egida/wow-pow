.PHONY: build

BUILD_ENVPARMS:=CGO_ENABLED=0

test:
	go test -race ./...

benchmark:
	go test -bench=. -benchtime=10s -benchmem ./...

docker-images:
	$(info Building docker image...)
	docker build --tag wow-pow-client -f deploy/client.Dockerfile .
	docker build --tag wow-pow-server -f deploy/server.Dockerfile .


docker-run-server:
	docker run -p 9000:9000  wow-pow-server:latest -- \
		--proof-token-size=40 \
		--proof-difficulty=24

docker-run-client:
	docker run --net host wow-pow-client:latest \
		--pow-concurrency=8 \
		--fetch-concurrency=1 \
		--pause-between-calls=2
