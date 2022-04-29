CREATE TABLE IF NOT EXISTS schedules
(
    id      SERIAL UNIQUE,
    chat_id BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS days
(
    id          SERIAL UNIQUE,
    schedule_id INT NOT NULL,
    caption     TEXT DEFAULT '',
    CONSTRAINT fk_schedule FOREIGN KEY (schedule_id)
        REFERENCES schedules (id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS pairs
(
    id          SERIAL UNIQUE,
    day_id      INT NOT NULL,
    title       TEXT DEFAULT '',
    information TEXT,
    CONSTRAINT fk_day FOREIGN KEY (day_id)
        REFERENCES days (id)
        ON DELETE CASCADE
);