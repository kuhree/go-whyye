.PHONY: clean build serve run dev

DATE=$(shell date +%Y%m%d)
SHORT_SHA=$(shell git rev-parse --short HEAD)
LONG_SHA=$(shell git rev-parse HEAD)

init:
	mkdir -p out out/bin out/share out/state
	go mod download

clean:
	rm -rf out/* 

run: init
	go run main.go

dev: init
	gow -c -v run .

build: init
	go build \
		-a -o ./out/bin/go-whyye \
		main.go

serve: build
	./out/bin/main

docker-pull:
	docker pull \
		git.littlevibe.net/kuhree/go-whyye:latest || echo "Failed to pull image."

docker-release: docker-pull docker-build
	docker push \
		git.littlevibe.net/kuhree/go-whyye

docker-dev: 
	docker run \
		-it --rm -e PORT=8080 -p 8080:8080 -e APP_ENV=development --entrypoint sh \
		git.littlevibe.net/kuhree/go-whyye:latest

docker-build:
	docker build \
		--tag go-whyye \
		--tag git.littlevibe.net/kuhree/go-whyye:latest \
		--tag git.littlevibe.net/kuhree/go-whyye:${DATE} \
		--tag git.littlevibe.net/kuhree/go-whyye:${SHORT_SHA} \
		--tag git.littlevibe.net/kuhree/go-whyye:${LONG_SHA} \
		--build-arg SENTRY_RELEASE=${LONG_SHA} \
		.

docker-serve: 
	docker run \
		-it --name go-whyye --rm -e PORT=8080 -p 8080:8080 \
		git.littlevibe.net/kuhree/go-whyye:latest
