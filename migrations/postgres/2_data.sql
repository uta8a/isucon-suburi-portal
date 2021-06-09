-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
INSERT INTO team (name) VALUES ('team1');
INSERT INTO team (name) VALUES ('team2');
INSERT INTO team (name) VALUES ('team3');

INSERT INTO score_log (team_id, score, message) VALUES (1, 3000, 'best score 1');
INSERT INTO score_log (team_id, score, message) VALUES (2, 200, 'best score 2');
INSERT INTO score_log (team_id, score, message) VALUES (3, 10, 'best score 3');

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE team;
DROP TABLE score_log;
