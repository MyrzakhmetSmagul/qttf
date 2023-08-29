-- Create the "city" table
CREATE TABLE city (
    city_id SERIAL PRIMARY KEY,
    city_name VARCHAR UNIQUE,
    city_link VARCHAR UNIQUE
);
-- Create the player table
CREATE TABLE player (
    player_id SERIAL PRIMARY KEY,
    player_name VARCHAR,
    player_surname VARCHAR,
    profile_link VARCHAR UNIQUE,
    city_id INT REFERENCES city(city_id) ON DELETE CASCADE
);
-- Create the "rating" table
CREATE TABLE rating (
    rating_id SERIAL PRIMARY KEY,
    player_id INT UNIQUE REFERENCES player(player_id) ON DELETE CASCADE,
    rating INT NOT NULL,
    last_update DATE NOT NULL
);