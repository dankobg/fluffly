# Fluffly (petfinder app)

Run caddy devproxy

1. `dc -p devproxy -f compose.devproxy.yaml up -d`

Fluffly local setup:

1. `docker compose up -d`

2. `just mg-up`

3. `just keto-create-tuples`

4. `go run main.go identities import-identities`

5. `go run main.go petfinder import-orgs --dir=~/Documents/petfinder_data --workers=32`

6. `go run main.go petfinder import-animals --dir=~/Documents/petfinder_data --workers=32`

7. `just devproxy-setup`

8. `just certs-trust`

Run fluffly:

1. `just dev`

2. `cd web && pnpm dev`
