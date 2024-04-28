-- +goose Up
-- +goose StatementBegin
ALTER TABLE "service_versions" ADD CONSTRAINT fk_service_versions_service FOREIGN KEY ("service_id") REFERENCES "services" ("service_id");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "service_versions" DROP CONSTRAINT fk_service_versions_service;
-- +goose StatementEnd
