APP=gogoSpoty

build:
	go build -o $(APP) ./cmd/$(APP)

run: build
	./$(APP)

docker:
	docker compose up --build

docker-down:
	docker compose down

release:
	GOOS=linux GOARCH=amd64 go build -o gogoSpoty-linux ./cmd/gogoSpoty/
	GOOS=darwin GOARCH=arm64 go build -o gogoSpoty-macos ./cmd/gogoSpoty/
	GOOS=windows GOARCH=amd64 go build -o gogoSpoty.exe ./cmd/gogoSpoty/

clean:
	rm -f $(APP)