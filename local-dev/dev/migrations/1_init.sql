DROP TABLE IF EXISTS team;
DROP TABLE IF EXISTS score_log;

CREATE TABLE team (
  id INT NOT NULL AUTO_INCREMENT,
  name TEXT NOT NULL,
  created_at datetime default current_timestamp,
  PRIMARY KEY (id)
);

CREATE TABLE score_log (
  id INT NOT NULL AUTO_INCREMENT,
  team_id INT NOT NULL,
  score INT,
  message TEXT,
  created_at datetime default current_timestamp,
  PRIMARY KEY (id)
);
