CREATE TABLE movements (
    movement_id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL
);

CREATE TABLE muscles (
    muscle_id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL
);

CREATE TABLE movement_muscles (
    movement_id INTEGER,
    muscle_id INTEGER,
    FOREIGN KEY (movement_id) REFERENCES movements(movement_id),
    FOREIGN KEY (muscle_id) REFERENCES muscles(muscle_id)
)
