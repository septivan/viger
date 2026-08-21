.PHONY: check test build up down

check:
	$(MAKE) -C backend check
	corepack pnpm --dir frontend check

test:
	$(MAKE) -C backend test
	corepack pnpm --dir frontend test

build:
	$(MAKE) -C backend build
	corepack pnpm --dir frontend build

up:
	docker compose up --build

down:
	docker compose down

