-- +goose Up
-- +goose StatementBegin
CREATE TABLE "service_versions" (
  "sv_id" SERIAL PRIMARY KEY,
  "version" VARCHAR(8),
  "changelog" TEXT,
  "service_id" INTEGER NOT NULL,
  "updated_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  "created_at" TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE "service_versions";
-- +goose StatementEnd
