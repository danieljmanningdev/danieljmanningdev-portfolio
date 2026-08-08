run:
	go run main.go

build:
	go build -o bin/app main.go

clean:
	rm -f app.db
	rm -rf bin/