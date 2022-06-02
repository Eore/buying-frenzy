-- restaurant
CREATE TABLE restaurants (
    id VARCHAR(32) PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    cash_balance DECIMAL NOT NULL
);

CREATE TABLE opening_hours (
    id VARCHAR(32) PRIMARY KEY,
    restaurant_id VARCHAR(32) NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    day VARCHAR(10) NOT NULL,
    start_time TIME NOT NULL,
    end_time TIME NOT NULL
);

CREATE TABLE menus (
    id VARCHAR(32) PRIMARY KEY,
    restaurant_id VARCHAR(32) NOT NULL REFERENCES restaurants(id) ON DELETE CASCADE,
    dish_name TEXT DEFAULT '',
    price DECIMAL NOT NULL
);


-- user
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    cash_balance DECIMAL NOT NULL
);

CREATE TABLE purchase_histories (
    id VARCHAR(32) PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    dish_name TEXT DEFAULT '',
    restaurant_name VARCHAR(100) NOT NULL,
    transaction_amount DECIMAL NOT NULL,
    transaction_date TIMESTAMP NOT NULL
);