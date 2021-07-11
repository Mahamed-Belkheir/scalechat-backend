CREATE TABLE messages (
  ID UUID,
  Room VARCHAR,
  Body TEXT,
  UserID UUID,
  CreatedAt BIGINT,
  PRIMARY KEY (ID),
);
