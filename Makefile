build:
	go build -o bin/app

run: build
	./bin/app

test:
	go test -v ./... -count=1

conn:
	psql "postgres://postgres:postgres@localhost:5432/gator"

up:
	cd sql/schema; \
	goose postgres postgres://postgres:postgres@localhost:5432/gator up

down:
	cd sql/schema; \
	goose postgres postgres://postgres:postgres@localhost:5432/gator down

