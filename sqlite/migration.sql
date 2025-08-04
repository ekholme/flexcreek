CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS movements (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    movement_type TEXT,
    movement_description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS workouts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    name TEXT,
    workout_date DATE NOT NULL DEFAULT CURRENT_DATE,
    notes TEXT,
    duration_seconds INTEGER,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS movement_instances (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    workout_id INTEGER NOT NULL,
    movement_id INTEGER NOT NULL,
    notes TEXT,
    rpe INTEGER, -- Rate of Perceived Exertion
    log_data TEXT, -- This will store our JSON blob
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (workout_id) REFERENCES workouts(id),
    FOREIGN KEY (movement_id) REFERENCES movements(id)
);


-- Indexes for performance

-- Note: The UNIQUE constraint on users(email) already creates an implicit index in SQLite,
-- so an explicit index on that column is not required for performance.

-- Indexes on foreign keys to speed up JOIN operations and filtering
CREATE INDEX idx_workouts_user_id ON workouts (user_id);
CREATE INDEX idx_movement_instances_workout_id ON movement_instances (workout_id);
CREATE INDEX idx_movement_instances_movement_id ON movement_instances (movement_id);

-- Index on other frequently queried columns
CREATE INDEX idx_workouts_workout_date ON workouts (workout_date);

-- Triggers to automatically update the 'updated_at' timestamp on row changes.
-- This ensures the timestamp is always current on any update, without
-- relying on the application layer to set it.

CREATE TRIGGER IF NOT EXISTS trigger_users_updated_at
AFTER UPDATE ON users
FOR EACH ROW
BEGIN
    UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;

CREATE TRIGGER IF NOT EXISTS trigger_movements_updated_at
AFTER UPDATE ON movements
FOR EACH ROW
BEGIN
    UPDATE movements SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;

CREATE TRIGGER IF NOT EXISTS trigger_workouts_updated_at
AFTER UPDATE ON workouts
FOR EACH ROW
BEGIN
    UPDATE workouts SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;

CREATE TRIGGER IF NOT EXISTS trigger_movement_instances_updated_at
AFTER UPDATE ON movement_instances
FOR EACH ROW
BEGIN
    UPDATE movement_instances SET updated_at = CURRENT_TIMESTAMP WHERE id = OLD.id;
END;