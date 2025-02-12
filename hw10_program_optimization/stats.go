package hw10programoptimization

import (
	"errors"
	"fmt"
	"io"
	"strings"

	//nolint:depguard
	jsoniter "github.com/json-iterator/go"
)

var jsoniterConfig = jsoniter.ConfigCompatibleWithStandardLibrary

type User struct {
	ID       int    `json:"Id"`       //nolint:tagliatelle
	Name     string `json:"Name"`     //nolint:tagliatelle
	Username string `json:"Username"` //nolint:tagliatelle
	Email    string `json:"Email"`    //nolint:tagliatelle
	Phone    string `json:"Phone"`    //nolint:tagliatelle
	Password string `json:"Password"` //nolint:tagliatelle
	Address  string `json:"Address"`  //nolint:tagliatelle
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
		if err := decoder.Decode(&user); errors.Is(err, io.EOF) {
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
