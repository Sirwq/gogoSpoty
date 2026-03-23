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
	GOOS=linux GOARCH=amd64 go build -tags standalone -o $(APP)-linux ./cmd/$(APP)/
	GOOS=darwin GOARCH=arm64 go build -tags standalone -o $(APP)-macos ./cmd/$(APP)/
	GOOS=windows GOARCH=amd64 go build -tags standalone -o $(APP).exe ./cmd/$(APP)/

clean:
	rm -f $(APP) $(APP)-linux $(APP)-macos $(APP).exe