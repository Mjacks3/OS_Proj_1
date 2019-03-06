DROP DATABASE IF EXISTS proj1;
CREATE DATABASE proj1;
USE proj1;

CREATE USER IF NOT EXISTS 'proj1user'@'localhost' IDENTIFIED BY 'password';
GRANT ALL PRIVILEGES ON * . * TO 'proj1user'@'localhost';
FLUSH PRIVILEGES;

CREATE TABLE sites (
    name varchar(128) PRIMARY KEY,
    role varchar(128) NOT NULL,
    uri varchar(256) NOT NULL
);

CREATE TABLE site_aps (
    ap varchar(128) PRIMARY KEY,
    url varchar(256) NOT NULL,
    name varchar(128) NOT NULL,
    FOREIGN KEY (name) REFERENCES sites(name) ON DELETE CASCADE
);
