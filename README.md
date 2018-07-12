# Serial UUID generator

Service for generating serial unique identifiers in UUID format. Can be used for systems with data shared across multiple databases.  

## How to use

### Build from source

Build and run go application. Setup connection to PostgreSQL database via environment variable `UUIDGEN_DATABASE_URL`.

```bash
dep ensure
go build
export UUIDGEN_DATABASE_URL=postgres://user:password@localhost/generator?sslmode=disable
./serial-uuid-generator
```

## Roadmap

* [x] storing sequence in PostgreSQL
* [x] generating sequential part by database sequence
* [x] request parameters validation
* [x] load config from environment
* [ ] logging
* [ ] database connection check on startup
* [ ] automatic table creation on startup
* [ ] statistics
* [ ] docker build in ci
* [ ] concurrency testing
* [ ] load testing
* [ ] quality and coverage badges
