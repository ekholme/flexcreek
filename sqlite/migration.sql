CREATE TABLE movements (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    movement_type TEXT,
    movement_description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE movement_instances (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    movement_id INTEGER NOT NULL,
    movement_sets INTEGER,
    reps_pers_set INTEGER,
    duration_seconds INTEGER,
    distance_meters INTEGER,
    rpe INTEGER,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (movement_id) REFERENCES movements (id) ON DELETE CASCADE
);

CREATE TABLE workouts (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    workout_date DATE NOT NULL DEFAULT CURRENT_DATE,
    notes TEXT,
    duration_seconds INTEGER
);

CREATE TABLE workout_movement_instances (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    workout_id INTEGER NOT NULL,
    movement_instance_id INTEGER NOT NULL,
    FOREIGN KEY (workout_id) REFERENCES workouts (id) ON DELETE CASCADE,
    FOREIGN KEY (movement_instance_id) REFERENCES movement_instances (id) ON DELETE CASCADE
);

CREATE TABLE users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    hashed_password TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP

);


--TODO CREATE INDEXES