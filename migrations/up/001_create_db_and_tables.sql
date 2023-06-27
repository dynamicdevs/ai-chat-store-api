\c mydb

CREATE EXTENSION IF NOT EXISTS vector WITH SCHEMA public;

CREATE TABLE product (
    id SERIAL PRIMARY KEY,
    sku TEXT NOT NULL,
    name TEXT NOT NULL,
    embedding vector(1536) NOT NULL,
    link TEXT,
    price TEXT
);

CREATE TABLE attribute (
    id SERIAL PRIMARY KEY,
    information TEXT NOT NULL,
    embedding vector(1536) NOT NULL
);

CREATE TABLE product_attribute (
    product_id INT NOT NULL REFERENCES product(id),
    attribute_id INT NOT NULL REFERENCES attribute(id),
    PRIMARY KEY (product_id, attribute_id)
);


CREATE TABLE client (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    system_prompt TEXT
);
