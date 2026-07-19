set dotenv-filename := ".env"
set dotenv-load

cwd := justfile_directory()

dev_mode := "nix" # nix or docker
domain := if env('DOMAIN', "") == "" { "fluffly-dev.xyz" } else { '-a "${DOMAIN}"' }
redis_host := if env('REDIS_HOST', "") == "" { "localhost" } else { '-a "${REDIS_HOST}"' }
redis_port := if env('REDIS_PORT', "") == "" { "6379" } else { '-a "${REDIS_PORT}"' }
redis_pwd := if env('REDIS_PASSWORD', "") == "" { "" } else { '-a "${REDIS_PASSWORD}"' }
db_uri := if env('DB_URI', "") == "" { "postgres://test:test@localhost:5432/test?sslmode=disable" } else { '-a "${DB_URI}"' }
migrations_dir := "db/migrations"

# db_host := if env('POSTGRES_HOST', "") == "" { "localhost" } else { '-a "${POSTGRES_HOST}"' }
# db_user := if env('POSTGRES_USER', "") == "" { "test" } else { '-a "${POSTGRES_USER}"' }
# db_password := if env('POSTGRES_PASSWORD', "") == "" { "test" } else { '-a "${POSTGRES_PASSWORD}"' }
# db_name := if env('POSTGRES_DB', "") == "" { "test" } else { '-a "${POSTGRES_DB}"' }
# db_sslmode := "disable"

# ----------------------------------------------------------------------------

default:
	@just -l

# Start server with live reload
dev:
  go tool -modfile=tools.mod air

# generate sql
gen-sql:
	@go generate ./...

# Generate openapi server
gen-openapi:
	go tool -modfile=tools.mod oapi-codegen --config=api/oapi-codegen.yaml api/schema/fluffly.yaml
	cd web && pnpm gen:openapi

# Run code generation
gen:
	go generate ./...
	just gen-openapi
	# just gen-sql

# Run go vet
vet:
	go vet ./...

# Run go fix
fix:
	go fix ./...

# Run go mod tidy
tidy:
	go mod tidy

# Show dependencies
deps:
	go mod graph

# List dependencies that have updates
# go list -u -m -f '{{if .Update}}{{.}}{{end}}' all
# go list -u -m -f '{{if (and (not (or .Main .Indirect)) .Update)}}{{.Path}}: {{.Version}} -> {{.Update.Version}}{{end}}' all
# deps-outdated:
# 	echo "TODO"

# Lint code with golangci-lint
lint:
	go tool -modfile=tools.mod golangci-lint run

# Lint every dockerfile in docker dir
lint-dockerfiles:
	docker container run -v {{cwd}}/docker:/dockerfiles --rm -i hadolint/hadolint hadolint --ignore DL3018 /dockerfiles/app.dockerfile /dockerfiles/devspace.dockerfile /dockerfiles/goreleaser.dockerfile

# Run go fmt
fmt:
	go fmt ./...

# Run all tests
test:
	go test ./...

# Run all tests with verbose flag
test-verbose:
	go test -v ./...

# Run all tests with race flag
test-race:
	go test -race ./...

# Run all tests with race and verbose flags
test-racev:
	go test -v -race ./...

# Run test coverage
test-coverage:
	go test -cover ./...

# ----------------------------------------------------------------------------

# Run docker system prune (without unused volumes)
[group('docker')]
prune:
	docker system prune -a -f --volumes

# Run docker volume prune (annonymous and unused)
[group('docker')]
prune-volumes:
	docker volume prune -a -f

# Run docker system prune all
[group('docker')]
prune-all: prune && prune-volumes

# ----------------------------------------------------------------------------

# Run docker compose commands
[group('docker')]
dc *flags:
	docker compose {{flags}}

# Run docker compose watch
[group('docker')]
dc-watch:
	docker compose watch

# Run docker compose watch without building & starting services
[group('docker')]
dc-watch-only:
	docker compose watch --no-up

# Run docker compose up -d
[group('docker')]
dc-up *flags:
	docker compose up -d {{flags}}

# Run docker compose down
[group('docker')]
dc-down:
	docker compose down

# Run docker compose start
[group('docker')]
dc-start:
	docker compose start

# Run docker compose stop
[group('docker')]
dc-stop:
	docker compose stop

# Run docker compose restart
[group('docker')]
dc-restart:
	docker compose restart

# ----------------------------------------------------------------------------

# Shell into a docker compose container by service name
[group('docker')]
sh name:
	docker compose exec -it {{name}} sh

# connect to redis via redis-cli
[group('docker')]
rediscli:
	docker compose exec -it redis redis-cli -h "{{redis_host}}" -p "{{redis_port}}" {{ if redis_pwd == "" { "" } else { redis_pwd } }}

# connect to postgres via psql
[group('docker')]
psql:
	docker compose exec -it pg psql "{{db_uri}}"

# ----------------------------------------------------------------------------

# Migrate the DB to the most recent version available
[group('migrate')]
mg-up *flags:
	go tool -modfile=tools.mod goose postgres "{{db_uri}}" -dir "{{migrations_dir}}" up {{flags}}

# Migrate the DB up by 1
[group('migrate')]
mg-up-by-one *flags:
	go tool -modfile=tools.mod goose postgres "{{db_uri}}" -dir "{{migrations_dir}}" up-by-one {{flags}}

# Migrate the DB to a specific VERSION
[group('migrate')]
mg-up-to version *flags:
	go tool -modfile=tools.mod goose postgres "{{db_uri}}" -dir "{{migrations_dir}}" up-to {{version}} {{flags}}

# Roll back the version by 1
[group('migrate')]
mg-down *flags:
	go tool -modfile=tools.mod goose postgres "{{db_uri}}" -dir "{{migrations_dir}}" down {{flags}}

# Roll back to a specific VERSION
[group('migrate')]
mg-down-to version *flags:
	go tool -modfile=tools.mod goose postgres "{{db_uri}}" -dir "{{migrations_dir}}" down-to {{version}} {{flags}}

# Re-run the latest migration
[group('migrate')]
mg-redo *flags:
	go tool -modfile=tools.mod goose postgres "{{db_uri}}" -dir "{{migrations_dir}}" redo {{flags}}

# Roll back all migrations
[group('migrate')]
mg-reset *flags:
	go tool -modfile=tools.mod goose postgres "{{db_uri}}" -dir "{{migrations_dir}}" reset {{flags}}

# Print the current version of the database
[group('migrate')]
mg-version *flags:
	go tool -modfile=tools.mod goose postgres "{{db_uri}}" -dir "{{migrations_dir}}" version {{flags}}

# Dump the migration status for the current DB
[group('migrate')]
mg-status *flags:
	go tool -modfile=tools.mod goose postgres "{{db_uri}}" -dir "{{migrations_dir}}" status {{flags}}

# migrations create new
[group('migrate')]
mg-create name *flags:
	go tool -modfile=tools.mod goose postgres "{{db_uri}}" -dir "{{migrations_dir}}" create {{name}} {{flags}}

# Apply sequential ordering to migrations
[group('migrate')]
mg-fix name *flags:
	go tool -modfile=tools.mod goose postgres "{{db_uri}}" -dir "{{migrations_dir}}" fix {{name}} {{flags}}

# Check migration files without running them
[group('migrate')]
mg-validate name *flags:
	go tool -modfile=tools.mod goose postgres "{{db_uri}}" -dir "{{migrations_dir}}" validate {{name}} {{flags}}

# migrate everything from scratch (use in dev only)
[group('migrate')]
mg-fresh:
	just pg-drop-schema
	just pg-create-schema
	just mg-up
	docker compose exec kratos /bin/sh -c "kratos migrate sql --read-from-env --config /etc/config/kratos/kratos.yaml --yes"
	docker compose exec keto /bin/sh -c "keto migrate up -c /etc/config/keto/keto.yaml --yes"

# seed initial data from users, orgs and animals
[group('migrate')]
mg-seed dir="~/Documents/petfinder_data":
  docker compose restart kratos keto
  just keto-create-tuples
  go run main.go identities import-identities
  go run main.go petfinder import-orgs --dir={{dir}} --workers=32
  go run main.go petfinder import-animals --dir={{dir}} --workers=32

# ----------------------------------------------------------------------------

# backup postgres with pg_dumpall
[group('backup')]
pg-backup:
	docker compose exec -it pg pg_dumpall -c -U test > {{cwd}}/db/dev_dump.sql

# restore postgres backup
[group('backup')]
pg-restore:
	docker compose exec -T pg psql "{{db_uri}}" < {{cwd}}/db/dev_dump.sql

# drop database
pg-drop-schema:
	docker compose exec -T pg psql "{{db_uri}}" -c "drop schema public cascade;"

# create database
pg-create-schema:
	docker compose exec -T pg psql "{{db_uri}}" -c "create schema public;"

# ----------------------------------------------------------------------------

# Import kratos identities
[group('ory')]
@kratos-import-identities:
	docker compose exec kratos /bin/sh -c "cd /etc/config/kratos/imports && kratos import identities developers.json customers.json -e http://localhost:4434 --format json"

# Create keto relation tuples
[group('ory')]
@keto-create-tuples:
	docker compose exec keto /bin/sh -c "keto relation-tuple create -f /etc/config/keto/relation-tuples -c /etc/config/keto/keto.yaml --format json --insecure-disable-transport-security"

# ----------------------------------------------------------------------------

# trust certs (too lazy to make it work on more distros, path is different)
certs-trust:
  #!/usr/bin/env sh
  sudo docker compose cp caddy:/data/caddy/pki/authorities/local/root.crt /etc/ca-certificates/trust-source/anchors/root.crt && \
  sudo chmod 644 /etc/ca-certificates/trust-source/anchors/root.crt && \
  sudo update-ca-trust

# ----------------------------------------------------------------------------
