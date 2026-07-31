CREATE TABLE games(
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    genre VARCHAR(50)
);

CREATE TABLE player_games (
    player_id INT REFERENCES players(id) ON DELETE CASCADE,
    game_id INT REFERENCES games(id) ON DELETE CASCADE,
    bought_at TIMESTAMP DEFAULT NOW(),

    PRIMARY KEY(player_id, game_id)
);