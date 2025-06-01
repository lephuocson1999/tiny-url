CREATE TABLE IF NOT EXISTS id_counters (
    name VARCHAR(32) PRIMARY KEY,
    value BIGINT NOT NULL
);

INSERT INTO id_counters (name, value)
VALUES ('url_id', 0)
ON CONFLICT (name) DO NOTHING; 