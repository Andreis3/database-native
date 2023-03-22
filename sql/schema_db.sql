CREATE SCHEMA goexpert;

-- Path: sql/schema_db.sql

CREATE TABLE goexpert.products (
    id varchar(255) NOT NULL PRIMARY KEY,
    name varchar(255) NOT NULL,
    price numeric(10,2) NOT NULL
);