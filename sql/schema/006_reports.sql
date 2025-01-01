-- +goose up
CREATE TABLE reports (
   report_id UUID PRIMARY KEY,
   created_at TIMESTAMP NOT NULL,
   updated_at TIMESTAMP NOT NULL,
   post_id UUID NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
   user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
   reason TEXT NOT NULL
);

-- +goose Down
DROP TABLE reports;
