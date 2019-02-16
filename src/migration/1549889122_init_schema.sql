-- +migrate Up
CREATE TABLE visits (
    id int UNSIGNED AUTO_INCREMENT,
    first_name varchar(255),
    last_name varchar(255),
    created_at timestamp default NOW(),
    updated_at timestamp default NOW(),
    PRIMARY KEY (id)
);

-- +migrate Down
DROP TABLE IF EXISTS visits;
