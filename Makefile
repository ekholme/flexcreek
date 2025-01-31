.PHONY: build-app run clean create-tables create-test-tables

#defining db variables
MIGRATION_SCRIPT := ./sqlite/migration.sql
DATABASE := flexcreek.db

#defining variables for the test db
MIGRATION_TEST_SCRIPT := ./sqlite/migration_test.sql
TEST_DATABASE := flexcreek_test.db

build-app:
	go build -o bin/flexcreek ./cmd/flexcreek.go

run: build-app
	@./bin/flexcreek

clean:
	@rm -rf bin

create-tables: $(MIGRATION_SCRIPT) $(DATABASE)
	sqlite3 $(DATABASE) < $(MIGRATION_SCRIPT)

create-test-tables: $(MIGRATION_TEST_SCRIPT) $(TEST_DATABASE)
	sqlite3 $(TEST_DATABASE) < $(MIGRATION_TEST_SCRIPT)