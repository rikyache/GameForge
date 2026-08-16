CREATE TABLE player_games(
     player_id INT NOT NULL,
     game_id INT NOT NULL,
     bought_at TIMESTAMP DEFAULT NOW(),

     PRIMARY KEY(player_id, game_id),

     FOREIGN KEY(player_id)
     REFERENCES players(id)
     ON DELETE CASCADE,

     FOREIGN KEY(game_id)
     REFERENCES games(id)
     ON DELETE CASCADE
);