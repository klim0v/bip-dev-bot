package service

import (
	"github.com/go-resty/resty/v2"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	pgql       *sqlx.DB
	httpClient *resty.Client
}

func NewRepository(pgql *sqlx.DB) *Repository {
	return &Repository{pgql: pgql, httpClient: resty.New()}
}

func (r *Repository) emailAddresses() []string {
	return []string{"klim0v-sergey@yandex.ru"}
}

func (r *Repository) minterAddresses() []string {
	return []string{"Mx00000000000000000000000000000987654321"}
}

func (r *Repository) btcAddresses() string {
	return "1K1AaFAChTdRRE2N4D6Xxz83MYtwFzmiPN"
}
