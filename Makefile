.PHONY: clean build serve run dev

init:
	mkdir -p out/bin out/share out/state

clean:
	rm -rf out/* 

build: init
	go build -o ./out/bin/main main.go

build-docker:
	docker build --tag go-whyye .

serve: build
	./out/bin/main

serve-docker: 
	docker run -it --rm -e PORT=8080 -p 8080:8080 go-whyye

run: init
	go run main.go

dev: init
	gow -c -v run .

