APP=gogoSpoty

build:
	go build -o $(APP) ./cmd/$(APP)

run: build
	./$(APP)

clean:
	rm -f $(APP)