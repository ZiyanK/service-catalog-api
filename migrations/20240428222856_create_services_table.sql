-- +goose Up
-- +goose StatementBegin
CREATE TABLE "services" (
  "service_id" SERIAL PRIMARY KEY,
  "name" VARCHAR(255) NOT NULL,
  "description" TEXT,
  "user_uuid" UUID NOT NULL,
  "updated_at" timestamp DEFAULT CURRENT_TIMESTAMP,
  "created_at" timestamp DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "services";
-- +goose StatementEnd
