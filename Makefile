# run: as running app
run:
	go run main.go

# test: run all test
test:
	go test -v ./...

# build: build docker image from Dockerfile
build:
	docker build -t usedeall:latest .

# run_container: run app in container docker
run_container: 
	docker run --name usedeall --network usedeall-network -p 8080:8080 usedeall:latest

# run_container_release: run app in container docker with gin mode release
run_container_release: 
	docker run --name usedeall --network usedeall-network -p 8080:8080 -e GIN_MODE=release usedeall:latest