CREATE TABLE IF NOT EXISTS `cars` (
    `id`           INTEGER PRIMARY KEY AUTOINCREMENT,
    `carName`      VARCHAR(64) NOT NULL UNIQUE
);
CREATE TABLE IF NOT EXISTS `rides` (
    `id`           INTEGER PRIMARY KEY AUTOINCREMENT,
    `driver`       VARCHAR(64) NOT NULL,
    `car`          INTEGER NOT NULL,
    `destination`  VARCHAR(64) NOT NULL,
    `created`      TEXT DEFAULT CURRENT_TIMESTAMP,
    `start`        TEXT NOT NULL,
    `end`          TEXT NOT NULL,
    FOREIGN KEY(car) REFERENCES cars(id)
);