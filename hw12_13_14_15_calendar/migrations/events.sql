-- +goose Up
CREATE TABLE events (
      id serial NOT NULL,
      title text NOT NULL,
      start_date text NOT NULL,
      end_date text NOT NULL,
      description text,
      user_id int NOT NULL,
      delay int,
      PRIMARY KEY(id)
);

-- +goose Down
DROP TABLE events;
