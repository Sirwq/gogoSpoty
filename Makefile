APP=gogoSpoty

build:
	go build -o $(APP) ./cmd/$(APP)

run: build
	./$(APP)

docker:
	docker compose up --build

docker-down:
	docker compose down

clean:
	rm -f $(APP)