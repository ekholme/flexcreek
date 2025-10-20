package sqlite

import (
	"database/sql"
	"testing"
)

// createAllTablesTestSchema creates the full database schema required for integration tests.
func createAllTablesTestSchema(t *testing.T, db *sql.DB) {
	t.Helper()
	// This schema is a simplified version of migration.sql, sufficient for testing.
	// It includes all tables because movement_instance depends on workouts and movements,
	// and workouts depends on users.
	schema := `
	CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		first_name TEXT, last_name TEXT, email TEXT UNIQUE, hashed_password TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS movements (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT UNIQUE, movement_type TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE TABLE IF NOT EXISTS workouts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER, workout_date DATE, notes TEXT, duration_seconds INTEGER,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
	);
	CREATE TABLE IF NOT EXISTS movement_instances (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		workout_id INTEGER, movement_id INTEGER, notes TEXT, rpe INTEGER, log_data TEXT,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (workout_id) REFERENCES workouts(id) ON DELETE CASCADE,
		FOREIGN KEY (movement_id) REFERENCES movements(id)
	);
	`
	if _, err := db.Exec(schema); err != nil {
		t.Fatalf("failed to create full schema: %v", err)
	}
}
