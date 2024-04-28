-- +goose Up
-- +goose StatementBegin
INSERT INTO users(user_uuid, email, password) VALUES ('d90f9b49-dcd9-4feb-8250-d013098e45ee', 'test@gmail.com', '$2a$04$YnrU9OJGE8ywjQ9yxIDZguyyRJucS4a5doOFk/3cde4OXxos6AkO.');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DELETE FROM users WHERE user_uuid = 'd90f9b49-dcd9-4feb-8250-d013098e45ee';
-- +goose StatementEnd
