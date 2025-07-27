package repository

import (
	"context"
	"fmt"
	"log"
	"rate_limiter/config"

	_ "github.com/jackc/pgx/v5/stdlib" // This line registers the "pgx" driver
	"github.com/jmoiron/sqlx"
)

type ListRepository struct {
	dsn string
	db  *sqlx.DB
	ctx context.Context
}

type IPRecord struct {
	IP string
}

func NewListRepository(ctx context.Context, config config.EnvConfig) (*ListRepository, error) {
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%d/%s",
		config.DB.User,
		config.DB.Pass,
		config.DB.Host,
		config.DB.Port,
		config.DB.Name,
	)
	repo := ListRepository{
		dsn: dsn,
		ctx: ctx,
	}
	err := repo.connect()
	if err != nil {
		return nil, err
	}

	return &repo, nil
}

func (r *ListRepository) GetWhitelist() ([]string, error) {
	return r.getList("select ip from whitelist")
}

func (r *ListRepository) GetBlacklist() ([]string, error) {
	return r.getList("select ip from blacklist")
}

func (r *ListRepository) getList(query string) ([]string, error) {
	rows, err := r.db.QueryxContext(r.ctx, query)
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

func (r *ListRepository) AddToWhitelist(network string) error {
	query := "insert into whitelist(ip) values ($1)"
	res, err := r.db.ExecContext(r.ctx, query, network)
	log.Printf("Add To WL DB: %+v", res)
	return err
}

func (r *ListRepository) AddToBlacklist(network string) error {
	query := "insert into blacklist(ip) values ($1)"
	_, err := r.db.ExecContext(r.ctx, query, network)
	return err
}

func (r *ListRepository) RemoveFromWhitelist(network string) error {
	query := "delete from whitelist where ip = $1"
	_, err := r.db.ExecContext(r.ctx, query, network)
	return err
}

func (r *ListRepository) RemoveFromBlacklist(network string) error {
	query := "delete from blacklist where ip = $1"
	_, err := r.db.ExecContext(r.ctx, query, network)
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
