CREATE TABLE IF NOT EXISTS `cars` (
    `id`           INTEGER PRIMARY KEY AUTOINCREMENT,
    `carName`      TEXT NOT NULL UNIQUE,
    `carColor`     TEXT
);
CREATE TABLE IF NOT EXISTS `rides` (
    `id`           INTEGER PRIMARY KEY AUTOINCREMENT,
    `driver`       TEXT NOT NULL,
    `car`          INTEGER NOT NULL,
    `destination`  TEXT NOT NULL,
    `created`      TEXT DEFAULT CURRENT_TIMESTAMP,
    `start`        TEXT NOT NULL,
    `end`          TEXT NOT NULL,
    FOREIGN KEY(car) REFERENCES cars(id)
);