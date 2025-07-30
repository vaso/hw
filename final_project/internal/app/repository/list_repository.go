package repository

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"rate_limiter/config"

	_ "github.com/jackc/pgx/v5/stdlib" // This line registers the "pgx" driver
)

type ListRepository struct {
	dsn string
	db  *sqlx.DB
}

type IPRecord struct {
	IP string
}

func NewListRepository(config config.EnvConfig) (*ListRepository, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		config.DB.User,
		config.DB.Pass,
		config.DB.Host,
		config.DB.Port,
		config.DB.Name,
	)
	repo := ListRepository{
		dsn: dsn,
	}
	err := repo.connect()
	if err != nil {
		return nil, err
	}

	return &repo, nil
}

func (r *ListRepository) GetWhitelist(ctx context.Context) ([]string, error) {
	return r.getList(ctx, "select ip from whitelist")
}

func (r *ListRepository) GetBlacklist(ctx context.Context) ([]string, error) {
	return r.getList(ctx, "select ip from blacklist")
}

func (r *ListRepository) getList(ctx context.Context, query string) ([]string, error) {
	rows, err := r.db.QueryxContext(ctx, query)
	if err != nil {
		return []string{}, err
	}
	defer rows.Close()

	var list []string
	for rows.Next() {
		var ip IPRecord
		err := rows.StructScan(&ip)
		if err != nil {
			return []string{}, err
		}
		list = append(list, ip.IP)
	}
	return list, nil
}

func (r *ListRepository) AddToWhitelist(ctx context.Context, network string) error {
	query := "insert into whitelist(ip) values ($1)"
	res, err := r.db.ExecContext(ctx, query, network)
	log.Printf("Add To WL DB: %+v", res)
	return err
}

func (r *ListRepository) AddToBlacklist(ctx context.Context, network string) error {
	query := "insert into blacklist(ip) values ($1)"
	_, err := r.db.ExecContext(ctx, query, network)
	return err
}

func (r *ListRepository) RemoveFromWhitelist(ctx context.Context, network string) error {
	query := "delete from whitelist where ip = $1"
	_, err := r.db.ExecContext(ctx, query, network)
	return err
}

func (r *ListRepository) RemoveFromBlacklist(ctx context.Context, network string) error {
	query := "delete from blacklist where ip = $1"
	_, err := r.db.ExecContext(ctx, query, network)
	return err
}

func (r *ListRepository) connect() error {
	db, err := sqlx.Open("pgx", r.dsn)
	if err != nil {
		return fmt.Errorf("connection error: %w", err)
	}
	r.db = db

	return nil
}

func (r *ListRepository) Close() error {
	return r.db.Close()
}
