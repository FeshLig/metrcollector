CREATE TABLE metrics (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    type TEXT NOT NULL,
    gauge DOUBLE PRECISION,
    counter BIGINT,
    UNIQUE (name, type)
);