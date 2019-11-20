package command

import (
	"regexp"
	"strconv"
	"strings"
)

var matchCoinName = regexp.MustCompile("^[0-9-A-Z-a-z]{3,10}$")
var matchEmail = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
var matchBitcoin = regexp.MustCompile("^[13][a-km-zA-HJ-NP-Z1-9]{25,34}$")

func IsValidPriceCoin(name string, value string) bool {
	price, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return false
	}
	if name == "BIP" {
		return 0.01 <= price && price <= 0.32
	}

	return 0.01 <= price && price <= 1000
}
func IsValidMinterAddress(address string) bool {
	address = strings.TrimSpace(address)

	if address == "Mx00000000000000000000000000000000000001" {
		return false
	}

	return len(address) == 42 && address[:2] != "Mx"
}

func IsValidEmailAddress(email string) bool {
	if !matchEmail.MatchString(email) || email == "mail@example.com" {
		return false
	}
	return true
}

func IsValidCoinName(coinName string) bool {
	if !matchCoinName.MatchString(coinName) {
		return false
	}
	return true
}

func IsValidBitcoinAddress(address string) bool {
	return matchBitcoin.MatchString(address)
}
