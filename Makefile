.PHONY: clean build serve run dev

include .env
include .env.production

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
		git.littlevibe.net/kuhree/go-whyye:latest  

docker-dev: 
	docker run \
		-it --rm -e PORT=8080 -p 8080:8080 -e APP_ENV=development --entrypoint sh \
		git.littlevibe.net/kuhree/go-whyye:latest

docker-build:
	docker build \
		--tag go-whyye --tag git.littlevibe.net/kuhree/go-whyye:latest \
		.

docker-serve: 
	docker run \
		-it --name go-whyye --rm -e PORT=8080 -p 8080:8080 \
		git.littlevibe.net/kuhree/go-whyye:latest
