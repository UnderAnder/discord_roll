.PHONY: clean

GOPATH = $(shell go env GOPATH)
GOBIN = $(GOPATH)/bin

### Custom Installs ########################################################
GO_MIGRATE = $(GOBIN)/migrate
$(GO_MIGRATE):
	@echo ">> Couldn't find go-migrate; installing..."
	go get -tags 'sqlite3' -u github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.0
	## Binary compiled without sqlite support
	#curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.0/migrate.linux-amd64.tar.gz | tar xvz
	#mv migrate $(GOBIN)/

### Targets ################################################################
all: clean build

MIGRATIONS = ./migrations
SQLITE_MIGRATIONS = $(MIGRATIONS)/sqlite
DATA = ./data
SQLITE_DB = $(DATA)/sqlite/bot.sqlite3
CITIES = $(DATA)/cities

migrate-up:
	@echo "Migrating up..."
	@mkdir -p ./data/sqlite
	@migrate -path $(SQLITE_MIGRATIONS) -database sqlite3://$(SQLITE_DB) up $(SCHEMA_VERSION)

migrate-down:
	@echo "Migrating down..."
	@migrate -path $(SQLITE_MIGRATIONS) -database sqlite3://$(SQLITE_DB) down $(SCHEMA_VERSION)

clean:
	@echo "Cleaning bin/..."
	@rm -rf bin/*

project-utils: $(GO_MIGRATE)
	@echo "Installing project utilities..."

extract-cities:
	@tar -zxf $(CITIES)/cities_ru.tar.gz

docker-image:
	@echo "Building docker image..."
	@docker build -t bot .

build:
	@echo "Building bot binary for use on local system..."
	@env CGO_ENABLED=1 go build -o ./bin/bot cmd/bot/main.go

run: build
	./bin/bot -t ${BOT_TOKEN}
