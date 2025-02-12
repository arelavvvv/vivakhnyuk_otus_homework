package hw10programoptimization

import (
	"fmt"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

var jsoniterConfig = jsoniter.ConfigCompatibleWithStandardLibrary

type User struct {
	ID       int    `json:"Id"`
	Name     string `json:"Name"`
	Username string `json:"Username"`
	Email    string `json:"Email"`
	Phone    string `json:"Phone"`
	Password string `json:"Password"`
	Address  string `json:"Address"`
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return countDomains(r, domain)
}

func countDomains(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	decoder := jsoniterConfig.NewDecoder(r)

	for {
		var user User
		if err := decoder.Decode(&user); err == io.EOF {
			break
		} else if err != nil {
			fmt.Printf("error decoding user: %v\n", err)
			continue
		}

		email := strings.ToLower(user.Email)
		if strings.HasSuffix(email, "."+domain) {
			parts := strings.SplitN(email, "@", 2)
			if len(parts) > 1 {
				domainName := parts[1]
				result[domainName]++
			}
		}
	}

	return result, nil
}
