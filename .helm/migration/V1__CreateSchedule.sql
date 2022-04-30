CREATE TABLE IF NOT EXISTS schedules
(
    id      UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chat_id BIGINT NOT NULL
);

CREATE TABLE IF NOT EXISTS days
(
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    schedule_id INT NOT NULL,
    caption     TEXT DEFAULT '',
    CONSTRAINT fk_schedule FOREIGN KEY (schedule_id)
        REFERENCES schedules (id)
        ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS pairs
(
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    day_id      INT NOT NULL,
    title       TEXT DEFAULT '',
    information TEXT,
    CONSTRAINT fk_day FOREIGN KEY (day_id)
        REFERENCES days (id)
        ON DELETE CASCADE
);