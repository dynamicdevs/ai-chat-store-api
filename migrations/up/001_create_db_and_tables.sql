
\c mydb

CREATE EXTENSION IF NOT EXISTS vector WITH SCHEMA public;


CREATE TABLE product (
    id SERIAL PRIMARY KEY,
    sku TEXT NOT NULL,
    name TEXT NOT NULL
);

CREATE TABLE attribute (
    id SERIAL PRIMARY KEY,
    product INT NOT NULL REFERENCES product(id),
    information TEXT NOT NULL,
    embedding vector(1536) NOT NULL
);
