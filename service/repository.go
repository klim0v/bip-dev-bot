package service

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/go-resty/resty/v2"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	rds        *redis.Client
	pgql       *sqlx.DB
	httpClient *resty.Client
}

func NewRepository(rds *redis.Client, pgql *sqlx.DB) *Repository {
	return &Repository{rds: rds, pgql: pgql, httpClient: resty.New()}
}

func (r *Repository) saveMinterAddressForBuy(chatID int64, minterAddressID string) error {
	if err := r.rds.Set(fmt.Sprintf("%d:minterAddress", chatID), minterAddressID, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (r *Repository) saveEmailAddressForBuy(chatID int64, emailAddressID string) error {
	if err := r.rds.Set(fmt.Sprintf("%d:emailAddress", chatID), emailAddressID, 0).Err(); err != nil {
		return err
	}
	return nil
}

type Address struct {
	ID    int
	Value string
}

func (r *Repository) emailAddresses() []Address {
	return []Address{
		{
			ID:    1,
			Value: "klim0v-sergey@yandex.ru",
		},
	}
}

func (r *Repository) minterAddresses() []Address {
	return []Address{
		{
			ID:    1,
			Value: "Mx00000000000000000000000000000987654321",
		},
	}
}

func (r *Repository) btcAddresses() string {
	return "1K1AaFAChTdRRE2N4D6Xxz83MYtwFzmiPN"
}
