CREATE TABLE metrics (
    id IDENTITY PRIMARY KEY,
    name VARCHAR(20) NOT NULL,
    type VARCHAR(10) NOT NULL,
    gauge DOUBLE PRECISION,
    counter BIGINT,
    UNIQUE (name, type)
);