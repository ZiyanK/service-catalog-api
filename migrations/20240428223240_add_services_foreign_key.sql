-- +goose Up
-- +goose StatementBegin
ALTER TABLE "services" ADD CONSTRAINT fk_services_users FOREIGN KEY ("user_uuid") REFERENCES "users" ("user_uuid");
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE "services" DROP CONSTRAINT fk_services_users;
-- +goose StatementEnd
