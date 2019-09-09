CREATE TABLE IF NOT EXISTS `cars`
(
    `id`       INTEGER PRIMARY KEY AUTOINCREMENT,
    `carName`  TEXT NOT NULL UNIQUE,
    `carColor` TEXT,
    `isBig`    INTEGER DEFAULT 0
);
CREATE TABLE IF NOT EXISTS `rides`
(
    `id`           INTEGER PRIMARY KEY AUTOINCREMENT,
    `driver`       TEXT NOT NULL,
    `car`          INTEGER,
    `destination`  TEXT NOT NULL,
    `created`      TEXT    DEFAULT CURRENT_TIMESTAMP,
    `start`        TEXT NOT NULL,
    `end`          TEXT NOT NULL,
    `confirmed`    INTEGER DEFAULT 0,
    `bigCarNeeded` INTEGER DEFAULT 0,
    FOREIGN KEY (car) REFERENCES cars (id)
);
CREATE TABLE IF NOT EXISTS `users`
(
    `id`       INTEGER PRIMARY KEY AUTOINCREMENT,
    `username` TEXT NOT NULL UNIQUE,
    `password` TEXT NOT NULL,
    `isAdmin`  INTEGER DEFAULT 0
);