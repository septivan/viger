.PHONY: check test build e2e up down

check:
	$(MAKE) -C backend check
	corepack pnpm --dir frontend check

test:
	$(MAKE) -C backend test
	corepack pnpm --dir frontend test

build:
	$(MAKE) -C backend build
	corepack pnpm --dir frontend build

e2e:
	corepack pnpm --dir frontend test:e2e

up:
	docker compose up --build

down:
	docker compose down
