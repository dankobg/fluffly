# Fluffly (petfinder app)

Project local setup:

1. `docker compose up -d`

2. `just mg-up`

3. `just keto-create-tuples`

4. `go run main.go identities import-identities`

5. `go run main.go petfinder import-orgs --dir=~/Documents/petfinder_data --workers=32`

6. `go run main.go petfinder import-animals --dir=~/Documents/petfinder_data --workers=32`

7. `just certs-trust`

Run project:

1. `just dev`

2. `cd web && pnpm dev`
