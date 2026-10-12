# Load the .env file
ifneq ($(wildcard .env),)
    include .env
    export $(shell sed 's/=.*//' .env)
endif

# Check if DB_URL is set
ifndef DB_URL
    $(error DB_URL is not set in .env file. Please create a .env file with DB_URL="your_connection_string")
endif

reset:
	sqlc generate
	cd sql/schema && goose postgres $(DB_URL) down && goose postgres $(DB_URL) up

reset-all:
	@read -p "This resets all migrations. Continue? [y/N] " answer; \
	case "$$answer" in \
		y|Y|yes|YES) \
			sqlc generate && cd sql/schema && goose postgres "$(DB_URL)" reset ;; \
		*) \
			echo "Cancelled." ;; \
	esac

up:
	sqlc generate
	cd sql/schema && goose postgres $(DB_URL) up

down:
	cd sql/schema && goose postgres $(DB_URL) down

build:
	go build -o Chirpy
