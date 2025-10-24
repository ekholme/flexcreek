CREATE TABLE
IF NOT EXISTS users
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    hashed_password BLOB NOT NULL,
    created_at        TEXT DEFAULT
(strftime
('%Y-%m-%d %H:%M:%S', 'now')),
    updated_at        TEXT DEFAULT
(strftime
('%Y-%m-%d %H:%M:%S', 'now'))
);

CREATE TABLE
IF NOT EXISTS activity_types
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE
IF NOT EXISTS workouts
(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL REFERENCES users
(id) ON
DELETE CASCADE,
    activity_type_id INTEGER
NOT NULL REFERENCES activity_types
(id),
    duration_minutes REAL,
    distance_miles REAL,
    workout_details TEXT,
    workout_date TEXT NOT NULL,
    created_at        TEXT DEFAULT
(strftime
('%Y-%m-%d %H:%M:%S', 'now')),
    updated_at        TEXT DEFAULT
(strftime
('%Y-%m-%d %H:%M:%S', 'now'))
);

-- create update triggers to automatically update the updated_at time
-- this ensures that I don't need to handle the updating in my Go code
CREATE TRIGGER update_workout_updated_at
AFTER
UPDATE ON workouts
FOR EACH ROW
WHEN NEW.updated_at = OLD.updated_at
-- Avoid recursive trigger calls
BEGIN
    UPDATE workouts SET updated_at = strftime('%Y-%m-%d %H:%M:%S', 'now') WHERE id = NEW.id;
END;

CREATE TRIGGER update_user_updated_at
AFTER
UPDATE ON users
FOR EACH ROW
WHEN NEW.updated_at = OLD.updated_at
-- Avoid recursive trigger calls
BEGIN
    UPDATE users SET updated_at = strftime('%Y-%m-%d %H:%M:%S', 'now') WHERE id = NEW.id;
END;

-- create indexes to speed up queries and joins
CREATE INDEX idx_workouts_user_id ON workouts (user_id);
CREATE INDEX idx_workouts_activity_type_id ON workouts (activity_type_id);
CREATE INDEX idx_workout_date ON workouts (workout_date);