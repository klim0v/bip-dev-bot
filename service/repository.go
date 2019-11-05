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

func (r *Repository) addMinterAddress(chatID int64, minterAddress string) (int, error) {
	return 1, nil
}

func (r *Repository) addBitcoinAddress(chatID int64, bitcoinAddress string) (int, error) {
	return 1, nil
}

func (r *Repository) savePriceForSell(chatID int64, price string) error {
	if err := r.rds.Set(fmt.Sprintf("%d:sell:price", chatID), price, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (r *Repository) saveMinterAddressForSell(chatID int64, price string) error {
	if err := r.rds.Set(fmt.Sprintf("%d:sell:price", chatID), price, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (r *Repository) saveCoinNameForSell(chatID int64, coinName string) error {
	if err := r.rds.Set(fmt.Sprintf("%d:sell:coinName", chatID), coinName, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (r *Repository) saveBitcoinAddressForSell(chatID int64, bitcoinAddressID int) error {
	if err := r.rds.Set(fmt.Sprintf("%d:sell:bitcoinAddressID", chatID), bitcoinAddressID, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (r *Repository) addEmailAddress(chatID int64, emailAddressID string) (int, error) {

	return 1, nil
}

func (r *Repository) saveEmailAddressForBuy(chatID int64, emailAddressID int) error {
	if err := r.rds.Set(fmt.Sprintf("%d:buy:emailAddressID", chatID), emailAddressID, 0).Err(); err != nil {
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

func (r *Repository) btcAddresses() []Address {
	return []Address{
		{
			ID:    1,
			Value: "1K1AaFAChTdRRE2N4D6Xxz83MYtwFzmiPN",
		},
	}
}
