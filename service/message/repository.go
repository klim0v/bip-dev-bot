package message

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/go-resty/resty/v2"
	"github.com/jmoiron/sqlx"
	"strings"
)

type Repository struct {
	rds        *redis.Client
	pgql       *sqlx.DB
	httpClient *resty.Client
}

func NewRepository(rds *redis.Client, pgql *sqlx.DB) *Repository {
	return &Repository{rds: rds, pgql: pgql, httpClient: resty.New()}
}

func (r *Repository) SellCoinName(chatID int64) (string, error) {
	name, err := r.rds.Get(keySellCoinName(chatID)).Result()
	if err != nil {
		return "", err
	}
	return name, nil
}

func keySellCoinName(chatID int64) string {
	return fmt.Sprintf("%d:sell:coinName", chatID)
}

//func (s *Application) SellCoinName(chatID int64) string {
//	name, err := s.sellCoinName(chatID)
//	if err != nil {
//		s.logger.Println(err)
//		return ""
//	}
//	return name
//}

func (r *Repository) AddMinterAddress(chatID int64, minterAddress string) (int, error) {
	return 1, nil
}

func (r *Repository) AddBitcoinAddress(chatID int64, bitcoinAddress string) (int, error) {
	return 1, nil
}

func (r *Repository) SaveSellPrice(chatID int64, price string) error {
	if err := r.rds.Set(fmt.Sprintf("%d:sell:price", chatID), price, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (r *Repository) SaveBuyMinterAddress(chatID int64, minterID int) error {
	if err := r.rds.Set(fmt.Sprintf("%d:buy:minter", chatID), minterID, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (r *Repository) SaveSellCoinName(chatID int64, coinName string) error {
	if err := r.rds.Set(fmt.Sprintf("%d:sell:coinName", chatID), strings.ToUpper(coinName), 0).Err(); err != nil {
		return err
	}
	return nil
}

func (r *Repository) SaveSellBitcoinAddress(chatID int64, bitcoinAddressID int) error {
	if err := r.rds.Set(fmt.Sprintf("%d:sell:bitcoinAddressID", chatID), bitcoinAddressID, 0).Err(); err != nil {
		return err
	}
	return nil
}

func (r *Repository) AddEmailAddress(chatID int64, emailAddressID string) (int, error) {

	return 1, nil
}

func (r *Repository) SaveEmailAddressForBuy(chatID int64, emailAddressID int) error {
	if err := r.rds.Set(fmt.Sprintf("%d:buy:emailAddressID", chatID), emailAddressID, 0).Err(); err != nil {
		return err
	}
	return nil
}

type Address struct {
	ID    int
	Value string
}

func (r *Repository) EmailAddresses() []Address {
	return []Address{
		{
			ID:    1,
			Value: "klim0v-sergey@yandex.ru",
		},
	}
}

func (r *Repository) MinterAddresses() []Address {
	return []Address{
		{
			ID:    1,
			Value: "Mx00000000000000000000000000000987654321",
		},
	}
}

func (r *Repository) BtcAddresses() []Address {
	return []Address{
		{
			ID:    1,
			Value: "1K1AaFAChTdRRE2N4D6Xxz83MYtwFzmiPN",
		},
	}
}
