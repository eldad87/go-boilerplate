-- +migrate Up
CREATE TABLE visits (
    id int,
    first_name varchar(255),
    last_name varchar(255),
    PRIMARY KEY (id)
);

-- +migrate Down
DROP TABLE IF EXISTS visits;
