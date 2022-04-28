CREATE TABLE IF NOT EXISTS schedules
(
    id   SERIAL UNIQUE,
    name VARCHAR(30)
);

CREATE TABLE IF NOT EXISTS days
(
    id          SERIAL UNIQUE,
    schedule_id INT NOT NULL,
    caption     TEXT DEFAULT '',
    CONSTRAINT fk_schedule FOREIGN KEY (schedule_id) REFERENCES schedules (id)
);

CREATE TABLE IF NOT EXISTS pairs
(
    id          SERIAL UNIQUE,
    day_id      INT NOT NULL,
    information TEXT,
    CONSTRAINT fk_day FOREIGN KEY (day_id) REFERENCES days (id)
);

INSERT INTO schedules
DEFAULT VALUES;

drop table if exists schedules cascade;
drop table if exists days cascade;
drop table if exists pairs cascade;
