-- +goose Up
-- +goose StatementBegin
CREATE INDEX idx_email ON users (email);
CREATE INDEX idx_user_uuid ON services USING HASH (user_uuid);
CREATE INDEX idx_sv_name_and_service_id ON service_versions (version, service_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP INDEX IF EXISTS idx_email;
DROP INDEX IF EXISTS idx_user_uuid;
DROP INDEX IF EXISTS idx_sv_name_and_service_id;
-- +goose StatementEnd
