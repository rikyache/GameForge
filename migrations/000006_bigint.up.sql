ALTER TABLE player_games
DROP CONSTRAINT player_games_player_id_fkey;

ALTER TABLE player_games
ALTER COLUMN player_id TYPE BIGINT;

ALTER TABLE players
ALTER COLUMN id TYPE BIGINT;

ALTER TABLE player_games
ADD CONSTRAINT player_games_player_id_fkey
FOREIGN KEY (player_id) REFERENCES players(id);