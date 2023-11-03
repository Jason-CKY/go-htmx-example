.DEFAULT_GOAL := help

# declares .PHONY which will run the make command even if a file of the same name exists
.PHONY: help
help:			## Help command
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)


.PHONY: lint
lint:			## Lint check
	docker run --rm -v $(PWD):/src:Z \
	--workdir=/src odinuge/yapf:latest yapf \
	--style '{based_on_style: pep8, dedent_closing_brackets: true, coalesce_brackets: true}' \
	--no-local-style --verbose --recursive --diff --parallel app

.PHONY: install-deps
install-deps:		## install all required node dependencies
	npm ci && cd src/app && npm ci

.PHONY: build-css
build-css:		## use tailwind cli to build out output css
	cd src/app && npx tailwindcss -i ./input.css -o ./static/output.css

.PHONY: watch-css
watch-css: 		## set css build into watch mode for development
	cd src.app && npx tailwindcss -i ./input.css -o ./static/output.css --watch

.PHONY: build-dev
local-dev:	local-dev	## starts golang application with live reload locally instead of with docker images
	cd src/app && air

.PHONY: build-dev
build-dev:	build-dev	## rebuild all the images in the docker-compose file
	docker-compose -f docker-compose.dev.yml up --build -d

.PHONY: start-dev
start-dev:		## deploy app in dev environment with hot reloading
	docker-compose -f docker-compose.dev.yml up -d

.PHONY: stop-dev
stop-dev:		## bring down all hosted services
	docker-compose -f docker-compose.dev.yml down

.PHONY: destroy-dev
destroy-dev:		## Bring down all hosted services with their volumes
	docker-compose -f docker-compose.dev.yml down -v

.PHONY: build
build:	build-css	## rebuild all the images in the docker-compose file
	docker-compose up --build -d

.PHONY: start
start:	build-css	## deploy app
	docker-compose up -d

.PHONY: stop
stop:		## bring down all hosted services
	docker-compose down

.PHONY: destroy
destroy:		## Bring down all hosted services with their volumes
	docker-compose down -v


.PHONY: build-tables
build-tables:		## initialize tables in directus for the app
	bash ./src/directus/build-tables.sh

.PHONY: seed-data
seed-data:		## seed initial data in directus for the app
	bash ./src/directus/seed-data.sh

.PHONY: set-permissions
set-permissions:		## set public access permissions for fastapi to perform CRUD unauthenticated
	bash ./src/directus/set-permissions.sh


.PHONY: initialize-db
initialize-db: build-tables seed-data set-permissions		## create service account, build-tables and seed data
	