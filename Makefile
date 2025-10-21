.PHONY: build-app run clean create-tables

#defining db variables
MIGRATION_SCRIPT := ./sqlite/migration.sql
POPULATE_ACTIVITIES_SCRIPT := ./sqlite/populate_activities.sql
DATABASE := flexcreek.db

build-app:
	go build -o bin/flexcreek ./cmd/flexcreek.go

run: build-app
	@./bin/flexcreek

clean:
	@rm -rf bin

create-tables: $(MIGRATION_SCRIPT) $(DATABASE)
	sqlite3 $(DATABASE) < $(MIGRATION_SCRIPT)

populate-activity-types-table: $(POPULATE_ACTIVITIES_SCRIPT) $(DATABASE)
	sqlite3 $(DATABASE) < $(POPULATE_ACTIVITIES_SCRIPT)

test:
	go test ./...