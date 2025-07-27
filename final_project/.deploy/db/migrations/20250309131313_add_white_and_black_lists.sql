-- +goose Up
-- +goose StatementBegin
INSERT INTO blacklist(ip) VALUES ('133.100.200.127/24'), ('230.100.100.1/32');
INSERT INTO blacklist(ip) VALUES ('133.200.200.127/24'), ('230.200.100.1/32');
INSERT INTO whitelist(ip) VALUES ('143.100.200.127/24'), ('240.100.100.1/32');
INSERT INTO whitelist(ip) VALUES ('143.200.200.127/24'), ('240.200.100.1/32');
-- +goose StatementEnd
-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
