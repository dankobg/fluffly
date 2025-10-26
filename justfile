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

sql_drop_public_tables := (
	"""
	DO $\\$
	DECLARE
		current_table text;
	BEGIN
		FOR current_table IN (SELECT table_name FROM information_schema.tables WHERE table_schema = 'public')
		LOOP
			EXECUTE 'DROP TABLE IF EXISTS public.' || current_table || ' CASCADE';
			END LOOP;
	END $\\$;
	"""
	)


# ----------------------------------------------------------------------------

default: 
	@just -l

# install needed tools
tools-install:
	go get -tool github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@latest
	go get -tool github.com/go-jet/jet/v2/cmd/jet@latest

# generate sql
gen-sql:
	@go generate ./...
	# jet -dsn={{db_uri}} -schema=public -path=./db/gen -ignore-tables continuity_containers,courier_message_dispatches,courier_messages,goose_db_version,identities,identity_credential_identifiers,identity_credential_types,identity_credentials,identity_login_codes,identity_recovery_addresses,identity_recovery_codes,identity_recovery_tokens,identity_registration_codes,identity_verifiable_addresses,identity_verification_codes,identity_verification_tokens,keto_relation_tuples,keto_uuid_mappings,networks,schema_migration,selfservice_errors,selfservice_login_flows,selfservice_recovery_flows,selfservice_registration_flows,selfservice_settings_flows,selfservice_verification_flows,session_devices,session_token_exchanges,sessions


# Generate openapi server
gen-openapi:
	oapi-codegen --config=api/oapi-codegen.yaml api/schema/fluffly.yaml
	#cd web && pnpm openapigen

# Run code generation
gen:
	go generate ./...
	just gen-openapi
	# just gen-sql

# Run go vet
vet:
	go vet ./...

# Run go mod tidy
tidy:
	go mod tidy

# Show dependencies
deps:
	go mod graph

# Lint code with golangci-lint
lint:
	golangci-lint run

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
prune:
	docker system prune -a -f --volumes

# Run docker volume prune (annonymous and unused)
prune-volumes:
	docker volume prune -a -f 

# Run docker system prune all
prune-all: prune && prune-volumes

# ----------------------------------------------------------------------------

# Run docker compose commands
dc *flags:
	docker compose {{flags}}

# Run docker compose watch
dc-watch:
	docker compose watch

# Run docker compose watch without building & starting services
dc-watch-only:
	docker compose watch --no-up

# Run docker compose up -d
dc-up *flags:
	docker compose up -d {{flags}}

# Run docker compose down
dc-down:
	docker compose down

# Run docker compose start
dc-start:
	docker compose start

# Run docker compose stop
dc-stop:
	docker compose stop

# Run docker compose restart
dc-restart:
	docker compose restart

# ----------------------------------------------------------------------------

# Shell into a docker compose container by service name
sh name:
	docker compose exec -it {{name}} sh

# connect to redis via redis-cli
rediscli:
	docker compose exec -it redis redis-cli -h "{{redis_host}}" -p "{{redis_port}}" {{ if redis_pwd == "" { "" } else { redis_pwd } }}
	
# connect to postgres via psql
psql:
	docker compose exec -it pg psql "{{db_uri}}"

# ----------------------------------------------------------------------------

# checks if goose is installed
_require_goose:
	#!/usr/bin/env sh
	command -v goose >/dev/null 2>&1 || { echo >&2 "goose is required, please install it to work with migrations"; exit 1; }

# Migrate the DB to the most recent version available
mg-up *flags: _require_goose
	goose postgres "{{db_uri}}" -dir "{{migrations_dir}}" up {{flags}}

# Migrate the DB up by 1
mg-up-by-one *flags: _require_goose
	goose postgres "{{db_uri}}" -dir "{{migrations_dir}}" up-by-one {{flags}}

# Migrate the DB to a specific VERSION
mg-up-to version *flags: _require_goose
	goose postgres "{{db_uri}}" -dir "{{migrations_dir}}" up-to {{version}} {{flags}}

# Roll back the version by 1
mg-down *flags: _require_goose
	goose postgres "{{db_uri}}" -dir "{{migrations_dir}}" down {{flags}}

# Roll back to a specific VERSION
mg-down-to version *flags: _require_goose
	goose postgres "{{db_uri}}" -dir "{{migrations_dir}}" down-to {{version}} {{flags}}

# Re-run the latest migration
mg-redo *flags: _require_goose
	goose postgres "{{db_uri}}" -dir "{{migrations_dir}}" redo {{flags}}

# Roll back all migrations
mg-reset *flags: _require_goose
	goose postgres "{{db_uri}}" -dir "{{migrations_dir}}" reset {{flags}}

# Print the current version of the database
mg-version *flags: _require_goose
	goose postgres "{{db_uri}}" -dir "{{migrations_dir}}" version {{flags}}
	
# Dump the migration status for the current DB
mg-status *flags: _require_goose
	goose postgres "{{db_uri}}" -dir "{{migrations_dir}}" status {{flags}}

# migrations create new
mg-create name *flags: _require_goose
	goose postgres "{{db_uri}}" -dir "{{migrations_dir}}" create {{name}} {{flags}}

# Apply sequential ordering to migrations
mg-fix name *flags: _require_goose
	goose postgres "{{db_uri}}" -dir "{{migrations_dir}}" fix {{name}} {{flags}}

# Check migration files without running them
mg-validate name *flags: _require_goose
	goose postgres "{{db_uri}}" -dir "{{migrations_dir}}" validate {{name}} {{flags}}

# migrate clean and up to latest
mg-fresh:
	just mg-reset
	just mg-up

# ----------------------------------------------------------------------------

# backup postgres with pg_dumpall
pg-backup:
	docker compose exec -it pg pg_dumpall -c -U test > {{cwd}}/db/dev_dump.sql

# restore postgres backup
pg-restore: pg-dropall
	docker compose exec -T pg psql "{{db_uri}}" < {{cwd}}/db/dev_dump.sql

# drop all tables in public schema
pg-dropall:
	docker compose exec -T pg psql "{{db_uri}}" -c "{{sql_drop_public_tables}}"

# ----------------------------------------------------------------------------

# Import kratos identities
kratos-import-identities:
	docker compose exec kratos /bin/sh -c "cd /etc/config/kratos/imports && kratos import identities developers.json customers.json -e http://localhost:4434 --format json"

# Create keto relation tuples
keto-create-tuples:
	docker compose exec keto /bin/sh -c "keto relation-tuple create /etc/config/keto/relation-tuples -c /etc/config/keto/keto.yaml --format json --insecure-disable-transport-security"

# ----------------------------------------------------------------------------

# checks if mkcert is installed
_require_mkcert:
	#!/usr/bin/env sh
	command -v mkcert >/dev/null 2>&1 || { echo >&2 "mkcert is required, please install it to work with certs"; exit 1; }

# install mkcert
certs-install: _require_mkcert
	mkcert -install

# uninstall mkcert
certs-uninstall: _require_mkcert
	mkcert -uninstall && rm -rf "$(mkcert -CAROOT)"

# generate certs
certs: _require_mkcert
	rm -f {{cwd}}/certs/local*.pem && \
	mkcert -cert-file /tmp/local-cert.pem -key-file /tmp/local-key.pem "{{domain}}" "*.{{domain}}" localhost 127.0.0.1 ::1 && \
	cp /tmp/local-{key,cert}.pem {{cwd}}/certs && \
	rm -f /tmp/local-{cert,key}.pem

