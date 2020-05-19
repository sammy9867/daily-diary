DROP DATABASE IF EXISTS diary_db;
CREATE DATABASE diary_db;
USE diary_db;

CREATE TABLE users (
    id BIGINT AUTO_INCREMENT, 
    username VARCHAR(255) UNIQUE, 
    email VARCHAR(100), 
    password VARCHAR(100),
    created_at DATETIME, 
    updated_at DATETIME,
	PRIMARY KEY (id)
);

CREATE TABLE entries (
    id BIGINT AUTO_INCREMENT, 
    title VARCHAR(255) NOT NULL,
    description VARCHAR(255) NOT NULL,
    owner_id BIGINT UNIQUE NOT NULL,
    created_at DATETIME, 
    updated_at DATETIME,
    PRIMARY KEY (id)
);

CREATE TABLE entry_images (
    id BIGINT AUTO_INCREMENT, 
    url VARCHAR(255) NOT NULL,
    entry_id BIGINT UNIQUE NOT NULL,
    created_at DATETIME, 
    updated_at DATETIME,
    PRIMARY KEY (id),
    FOREIGN KEY (entry_id) REFERENCES entries(id)
);