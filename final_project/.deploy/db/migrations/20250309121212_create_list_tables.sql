-- +goose Up
-- +goose StatementBegin
CREATE TABLE blacklist (
    ip varchar(20)
);
CREATE TABLE whitelist (
    ip varchar(20)
);
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
DROP TABLE blacklist;
DROP TABLE whitelist;
-- +goose StatementEnd

