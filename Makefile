.PHONY: clean build serve run dev

init:
	mkdir -p out out/bin out/share out/state

clean:
	rm -rf out/* 

run: init
	go run main.go

dev: init
	gow -c -v run .

build: init
	go build -o ./out/bin/main main.go

serve: build
	./out/bin/main

docker-pull:
	docker pull git.littlevibe.net/kuhree/go-whyye:latest  

docker-push: docker-pull docker-build
	docker push git.littlevibe.net/kuhree/go-whyye:latest  

docker-dev: 
	docker run -it --rm -e PORT=8080 -p 8080:8080 -e APP_ENV=development --entrypoint bash go-whyye

docker-build:
	docker build --tag git.littlevibe.net/kuhree/go-whyye:latest .

docker-serve: 
	docker run -it --rm -e PORT=8080 -p 8080:8080 go-whyye

