-- +goose Up
-- +goose StatementBegin
CREATE TABLE "users" (
  "user_uuid" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "email" VARCHAR(50) UNIQUE NOT NULL,
  "password" VARCHAR(255) NOT NULL,
  "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "users";
-- +goose StatementEnd
