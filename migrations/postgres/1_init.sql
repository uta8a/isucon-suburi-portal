-- +migrate Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TABLE team (
  id serial primary key,
  name text not null,
  created_at timestamp not null default current_timestamp
);

CREATE TABLE score_log (
  id serial primary key,
  team_id int references team(id),
  score int,
  message text,
  created_at timestamp not null default current_timestamp
);

-- +migrate Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TABLE team;
DROP TABLE score_log;
