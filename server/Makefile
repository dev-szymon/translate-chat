build:
	go build -o ./bin/server ./cmd/main.go 
start: build
	./bin/server
docker_build: 
	docker build -t translate-chat-server .
docker_run: 
	docker run -p 8055:8055 --name translate-chat-server translate-chat-server